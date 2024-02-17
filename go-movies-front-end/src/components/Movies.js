import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

const Movies = () => {
    const [movies, setMovies] = useState([]);

    useEffect( () => {
        /* let moviesList =[
            {
                id: 1,
                title: "Highlander",
                release_date: "1986-03-07",
                runtime: 116,
                mpaa_rating: "R",
                description: "Some long description",
            },
            {
                id: 1,
                title: "Riders",
                release_date: "1990-03-15",
                runtime: 115,
                mpaa_rating: "PG-13",
                description: "Some long description",
            },
        ];

        setMovies(moviesList); */
        console.log("Use Effect called")
        const headers = new Headers();
        headers.append("Content-Type","application/json");
        const requestOptions = {
            method: "GET",
            headers: headers
        }
    
        fetch(`http://localhost:8080/movies`,requestOptions)
            .then( (response) => {
                console.log(response)
                return response.json()
            })
            .then( (data) => {
                console.log(data)
                setMovies(data)
            })
            .catch( (err => {
                console.log(err)
            }))

    }, []);

    return(
        <>
        <div>
            <h2>Movies</h2>
            <hr></hr>
            <table className="table table-striped table-hover">
                <thead>
                    <tr>
                        <th>Movie</th>
                        <th>Release Date</th>
                        <th>Rating</th>
                    </tr>
                </thead>
                <tbody>
                    {movies.map((m) => (
                        <tr key={m.id}>
                            <td>
                                <Link to={`/movies/${m.id}`}>
                                    {m.title}
                                </Link>
                            </td>
                            <td>{m.release_date}</td>
                            <td>{m.mpaa_rating}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
        </>
    )
}

export default Movies;