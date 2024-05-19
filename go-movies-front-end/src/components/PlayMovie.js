import { Player } from 'video-react';
import { useState, useEffect } from 'react';

const PlayMovie = (props) => {
    const [videoUrl, setVideoUrl] = useState('');
    const [errorMsg, setErrorMsg] = useState('');

    useEffect(() => {
        const fetchVideo = async () => {
            try {
            const headers = new Headers();
            headers.append("Content-Type", "application/json")

            const requestOptions = {
                method: "GET",
                headers: headers
                }
              const response = await fetch(`/movie/${props.movieID}/video`,requestOptions);
              console.log(response.status);
              if (response.ok) {
                //console.log(response.status);
                const blob = await response.blob();
                const url = URL.createObjectURL(blob);
                //console.log(url);
                setVideoUrl(url);
              } else {
                setErrorMsg("Unfortunately, Movie video is not availble at the moment. Pls try later.");
              }
              
            } catch (error) {
              console.error('Error fetching video:', error);
            }
          };
          fetchVideo();
    }, [props.movieID]);
    return (
        <div>
            {videoUrl !== '' ? (
                    <Player
                        playsInline
                        src={videoUrl}
                        fluid={true}
                        width={640}
                        height={360}
                    />
                ) : (
                    <p>{errorMsg}</p>
                )
            }
        </div>
    )

}

export default PlayMovie;