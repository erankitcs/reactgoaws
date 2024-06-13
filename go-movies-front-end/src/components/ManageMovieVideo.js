import { useLocation, useParams, useOutletContext } from "react-router-dom";
import { useState, useEffect, useRef } from "react";
import Input from "./form/Input";
import SpinnerButton from "./form/SpinnerButton";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrashCan, faCheck, faCheckSquare } from '@fortawesome/free-solid-svg-icons';
import Swal from 'sweetalert2'
import DateFormater from "../utils/DateFormater";

const ManageMovieVideo = () => {
    //const navigate = useNavigate();
    const location = useLocation();
    const { movieTitle } = location.state;
    let { id } = useParams();
    const { jwtToken } = useOutletContext();
    const [ movieVideoFile, setMovieVideoFile ] = useState(null);
    const [ movieVideos, setMovieVideos ] = useState([]);
    const [ movieVideo, setMovieVideo ] = useState({});
    const [isDeleting, setIsDeleting] = useState({});
    const [isMarkingLatest, setIsMarkingLatest] = useState({});
    
    //const [ error, setError ] = useState(null);
    const [ errors, setErrors ] = useState([]);

    // Var to reference input file element
    const inputVideoFile = useRef(null);

    const Toast = Swal.mixin({
        toast: true,
        position: "top-end",
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.onmouseenter = Swal.stopTimer;
          toast.onmouseleave = Swal.resumeTimer;
        }
      });

    const hasError = (key) => {
        return errors.indexOf(key) !== -1;
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
                // Creating a list of deleting and marking latest flags with key as video id
                let deleting = {};
                let markingLatest = {};
                data.forEach(video => {
                    deleting[video.id] = false;
                    markingLatest[video.id] = false;
                })
                // Setting the flags
                setIsDeleting(deleting);
                setIsMarkingLatest(markingLatest);
            }
        })
        .catch( err => {
            console.log(err)
        })

    }, [id, jwtToken])
    

    const handleChange = () => (event) =>{
        //let value = event.target.value;
        //let name = event.target.name;
        let file = event.target.files[0];
        setMovieVideoFile(file)
        setErrors([]);
        
    }

    const handleSubmit = (event) => {
        console.log("submit");
        event.preventDefault();
        let errors = []
        console.log(movieVideoFile);
        if (movieVideoFile === null) {
            console.log("file not found");
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

        fetch(`/admin/movies/${id}/video`, requestOptions)
            .then((repsonse) => repsonse.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                    Toast.fire({
                        icon: "error",
                        title: "Failed to upload the video"
                      });
                } else {
                    Toast.fire({
                        icon: "success",
                        title: "Video has been uploaded"
                      });
                    // Mark old video as not latest
                    let updatedMovieVideos = movieVideos.map(video => {
                            video.is_latest = false;
                        return video;
                    });
                    // Add the new video to the list of videos
                    updatedMovieVideos = [...movieVideos, data];
                    setMovieVideos(updatedMovieVideos);
                    setMovieVideos(updatedMovieVideos);
                    // Reset the form
                    setMovieVideoFile(null);
                    setMovieVideo({});
                    inputVideoFile.current.value = "";
                    //navigate("/manage-catalogue")
                }
            })
            .catch( err => {
                console.log(err)
            })

    }

    const  handleDelete = (videoID) => {
        console.log("delete-", id, videoID);
        Swal.fire({
            title: "Delete Movie Video?",
            text: "You won't be able to undo this!",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Yes, delete it!"
          }).then((result) => {
            if (result.isConfirmed) {
                setIsDeleting((prevState) => ({
                    ...prevState,
                    [videoID]: true,
                  }));
                const headers = new Headers()
                headers.append("Content-Type", "application/json");
                headers.append("Authorization", "Bearer "+ jwtToken);
        
                let requestOptions = {
                    method: "DELETE",
                    headers: headers,
                    credentials: "include"
                }
        
                fetch(`/admin/movies/${id}/videos/${videoID}`, requestOptions)
                    .then((repsonse) => repsonse.json())
                    .then((data) => {
                        if (data.error) {
                            console.log(data.error)
                            Toast.fire({
                                icon: "error",
                                title: "Failed to delete the video"
                              });
                        } else {
                            Toast.fire({
                                icon: "success",
                                title: "Video has been deleted"
                              });
                              // Remove this movie from movieVideos
                              let updatedMovieVideos = movieVideos.filter(video => video.id !== videoID);
                              setMovieVideos(updatedMovieVideos);
                        }
                    })
                    .catch( err => {
                        console.log(err)
                        Toast.fire({
                            icon: "error",
                            title: "Failed to delete the video"
                          });
                    })
                  
                setIsDeleting((prevState) => ({
                    ...prevState,
                    [videoID]: false,
                  }));
            }
          })
    }

    const  handleMarkLatest = (videoID, vidoePath) => {
        console.log("Mark latest",id, videoID);
        setIsMarkingLatest((prevState) => ({
            ...prevState,
            [videoID]: true,
          }));
         // console.log(isDeleting);
          // Wait for 1 min
          // Wait for 1 minute
            const headers = new Headers()
            headers.append("Content-Type", "application/json");
            headers.append("Authorization", "Bearer "+ jwtToken);

            setMovieVideo({
                movie_id: parseInt(id),
                id: parseInt(videoID),
                video_path: vidoePath,
                is_latest: true
            })

            let requestOptions = {
                body: JSON.stringify(movieVideo),
                method: "PATCH",
                headers: headers,
                credentials: "include"
            }

            fetch(`/admin/movies/${id}/videos/${videoID}`, requestOptions)
            .then((repsonse) => repsonse.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                    Toast.fire({
                        icon: "error",
                        title: "Failed to mark the video as latest"
                      });
                } else {
                    Toast.fire({
                        icon: "success",
                        title: "Video has been marked as latest."
                      });
                    // update each movie video to is_latest false but mark current movie as latest
                    let updatedMovieVideos = movieVideos.map(video => {
                        if (video.id === videoID) {
                            video.is_latest = true;
                        } else {
                            video.is_latest = false;
                        }
                        return video;
                    });
                    setMovieVideos(updatedMovieVideos);
                }
            })
            .catch( err => {
                console.log(err)
                Toast.fire({
                    icon: "error",
                    title: "Failed to mark the video as latest"
                  });
            })
          
        setIsMarkingLatest((prevState) => ({
            ...prevState,
            [videoID]: false,
          }));

    }

    

    return(
        <div>
            <h2>Movie: { movieTitle }</h2>
            <hr></hr>
            <form onSubmit={handleSubmit} encType="multipart/form-data">
            <Input
                    title={"Select a Video File"}
                    className={"form-control"}
                    type={"file"}
                    name={"movie_video_file"}
                    //value={movieVideoFile}
                    onChange={handleChange("movie_video_file")}
                    errorDiv={ hasError("movie_video_file") ? "text-danger": "d-none" }
                    errorMsg={"Please select a video"}
                    ref={inputVideoFile}
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
                    <th>Actions</th>
                </tr>
            </thead>
        <tbody>
            {movieVideos.map((mv) => (
                <tr key={mv.id}>
                    <td>{mv.video_path}</td>
                    <td>{DateFormater(mv.created_at)}</td>
                    <td>{mv.is_latest && <button type="button" class="btn btn-outline-success btn-sm me-1" disabled> <FontAwesomeIcon icon={faCheckSquare} /></button>}
                    </td>
                    <td>
                        <SpinnerButton
                            className="btn btn-outline-danger btn-sm me-1"
                            onClick={() => handleDelete(mv.id)}
                            faIcon={faTrashCan}
                            // set state isDeleting with key as video id
                            state={isDeleting[mv.id]}
                        />
                        {!mv.is_latest && <SpinnerButton
                            className="btn btn-outline-success btn-sm me-1"
                            onClick={() => handleMarkLatest(mv.id, mv.video_path)}
                            faIcon={faCheck}
                            state={isMarkingLatest[mv.id]}
                        />
                        }
                    </td>
                </tr>
            ))}
        </tbody>
    </table>
</div>
    )
}

export default ManageMovieVideo;