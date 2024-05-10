import { useLocation, useNavigate, useParams, useOutletContext } from "react-router-dom";
import { useState, useEffect } from "react";
import Input from "./form/Input";
const UploadMovieVideo = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { movieTitle } = location.state;
    let { id } = useParams();
    const { jwtToken } = useOutletContext();
    const [ movieVideoFile, setMovieVideoFile ] = useState(null);
    const [ movieVideos, setMovieVideos ] = useState([]);
    
    //const [ error, setError ] = useState(null);
    const [ errors, setErrors ] = useState([]);

    const hasError = (key) => {
        return errors.indexOf(key) !== -1;
    }

    const formatDate = (date) => {
        const d = new Date(date);
        return d.toUTCString();
    }

    useEffect(() => {
      // Getting all the video history for given movie Id
        const headers = new Headers()
        headers.append("Content-Type", "application/json");
        headers.append("Authorization", "Bearer "+ jwtToken);

        let requestOptions = {
            method: "GET",
            headers: headers,
            credentials: "include"
        }

        fetch(`/admin/movies/${id}/videos`, requestOptions)
        .then((repsonse) => repsonse.json())
        .then((data) => {
            if (data.error) {
                console.log(data.error)
            } else {
                setMovieVideos(data)
            }
        })
        .catch( err => {
            console.log(err)
        })

    }, [id, jwtToken])
    

    const handleChange = () => (event) =>{
        let value = event.target.value;
        let name = event.target.name;
        let file = event.target.files[0];
        console.log(value);
        console.log(name);
        console.log(file);
        setMovieVideoFile(file)
        
    }

    const handleSubmit = (event) => {
        console.log("submit");
        event.preventDefault();
        let errors = []
        console.log(movieVideoFile);
        if (movieVideoFile === null) {
            console.log("file found");
            errors.push("movie_video_file");
            setErrors(errors);
            return false;
        }

        // validation completed proceed with video upload
        // Calculate the number of chunks and the size of each chunk
        const fileSize = movieVideoFile.size;
        const chunkSize = 1024 * 1024; // Set chunk size to 1MB
        const totalChunks = Math.ceil(fileSize / chunkSize);
        var formData = new FormData();        
        formData.append('movievideofile', movieVideoFile);
        formData.append('totalChunks', totalChunks);
        //console.log(formData);
        const headers = new Headers()
        //headers.append("Content-Type", "application/json");
        headers.append("Authorization", "Bearer "+ jwtToken);

        let requestOptions = {
            body: formData,
            method: "POST",
            headers: headers,
            credentials: "include"
        }

        fetch(`/admin/movies/${id}/upload`, requestOptions)
            .then((repsonse) => repsonse.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                } else {
                    navigate("/manage-catalogue")
                }
            })
            .catch( err => {
                console.log(err)
            })

    }

    

    return(
        <div>
            <h2>Movie: { movieTitle }</h2>
            <hr></hr>
            <form onSubmit={handleSubmit} enctype="multipart/form-data">
            <Input
                    title={"Select a Video File"}
                    className={"form-control"}
                    type={"file"}
                    name={"movie_video_file"}
                    //value={movieVideoFile}
                    onChange={handleChange("movie_video_file")}
                    errorDiv={ hasError("movie_video_file") ? "text-danger": "d-none" }
                    errorMsg={"Please select a video"}
            />
            <button className="btn btn-primary">Upload</button>
        </form>
        
        <h2 className="mt-4">Upload history</h2>
        <hr></hr>
        <table className="table table-striped table-hover">
            <thead>
                <tr>
                    <th>VidoePath</th>
                    <th>Upload Time</th>
                    <th>Latest</th>
                </tr>
            </thead>
        <tbody>
            {movieVideos.map((mv) => (
                <tr key={mv.id}>
                    <td>{mv.video_path}</td>
                    <td>{formatDate(mv.created_at)}</td>
                    <td>{mv.is_latest.toString()}</td>
                </tr>
            ))}
        </tbody>
    </table>
</div>
    )
}

export default UploadMovieVideo;