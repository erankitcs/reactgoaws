import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import PlayMovie from "./PlayMovie";
import MovieChat from "./MovieChat";
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
 

const Movie = () => {
    const [movie, setMovie] = useState({});
    let { id } = useParams();
    useEffect( () => {
        // let myMovie = {
        //         id: 1,
        //         title: "Highlander",
        //         release_date: "1986-03-07",
        //         runtime: 116,
        //         mpaa_rating: "R",
        //         description: "Some long description",
        // }
        // setMovie(myMovie)
        const headers = new Headers();
        headers.append("Content-Type", "application/json")

        const requestOptions = {
            method: "GET",
            headers: headers
        }

        fetch(`/movie/${id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setMovie(data)
            })
            .catch(err => {
                console.log(err)
            })
    }, [id]);

    if (movie.genres) {
        movie.genres = Object.values(movie.genres)
    } else {
        movie.genres =[];
    }

    return(
        <Container>
        <Row>
         <Col sm={10}>
            <h2>Movie: {movie.title}</h2>
            <small><em>{movie.release_date}, {movie.runtime} minutes, Rated {movie.mpaa_rating}</em></small>
            <br></br>
            {movie.genres.map( (g) => 
                <span key={g.genre} className="badge bg-secondary me-2">{g.genre}</span>
            )}
            <hr></hr>
            {
                // movie.image !== ""  && 
                // <div className="mb-3">
                //     <img src={`https://image.tmdb.org/t/p/w200/${movie.image}`} alt="poster" />
                // </div>
            }
            <p>{movie.description}</p>
            <hr></hr>
         </Col>
        </Row>
        
        <Row>
                <Col sm={10}>
                    <PlayMovie
                        movieTitle={movie.title}
                        movieID={id}
                    ></PlayMovie>
                </Col>
        </Row>
        <Row >
                <Col sm={10}>
                <hr></hr>
                 <MovieChat></MovieChat>
                </Col>
        </Row>

            <hr></hr>
   
        </Container>
    )
}

export default Movie;