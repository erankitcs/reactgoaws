package graph

import (
	"backend/internal/models"
	"errors"
	"strings"

	"github.com/graphql-go/graphql"
)

type Graph struct {
	Movies      []*models.Movie
	QueryString string
	Config      graphql.SchemaConfig
	fields      graphql.Fields
	movieType   *graphql.Object
}

func New(movies []*models.Movie) *Graph {

	var movieType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Movie",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"description": &graphql.Field{
					Type: graphql.String,
				},
				"year": &graphql.Field{
					Type: graphql.Int,
				},
				"release_date": &graphql.Field{
					Type: graphql.DateTime,
				},
				"runtime": &graphql.Field{
					Type: graphql.Int,
				},
				"rating": &graphql.Field{
					Type:        graphql.Int,
					Description: "1-5",
				},
				"mpaa_rating": &graphql.Field{
					Type:        graphql.String,
					Description: "PG-13, R, etc",
				},
				"created_at": &graphql.Field{
					Type:        graphql.DateTime,
					Description: "When the movie was created",
				},
				"updated_at": &graphql.Field{
					Type:        graphql.DateTime,
					Description: "When the movie was updated",
				},
				"image": &graphql.Field{
					Type:        graphql.String,
					Description: "Image url",
				},
			},
		},
	)

	var fields = graphql.Fields{
		"list": &graphql.Field{
			Type:        graphql.NewList(movieType),
			Description: "Get all movies",
			Resolve:     func(p graphql.ResolveParams) (interface{}, error) { return movies, nil },
		},
		"search": &graphql.Field{
			Type:        movieType,
			Description: "Search movies by title",
			Args: graphql.FieldConfigArgument{
				"titleContains": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// get the movies containing title as passed argument
				var theMovieList []*models.Movie
				search, ok := p.Args["titleContains"].(string)
				if ok {
					for _, currentMovie := range movies {
						// check if the movie title contains the search string with upper and lower case search
						// and add it to the list
						if strings.Contains(strings.ToLower(currentMovie.Title), strings.ToLower(search)) {
							theMovieList = append(theMovieList, currentMovie)
						}
					}
				}

				return theMovieList, nil
			},
		},
		"get": &graphql.Field{
			Type:        movieType,
			Description: "Get a movie by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, movie := range movies {
						if movie.ID == id {
							return movie, nil
						}
					}
				}
				return nil, nil
			},
		},
	}

	return &Graph{
		Movies:    movies,
		fields:    fields,
		movieType: movieType,
	}
}

// Function to query graphql
func (g *Graph) Query() (*graphql.Result, error) {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: g.fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	params := graphql.Params{Schema: schema, RequestString: g.QueryString}
	resp := graphql.Do(params)
	if len(resp.Errors) > 0 {
		return nil, errors.New("error in executing the query")
	}
	return resp, nil
}
