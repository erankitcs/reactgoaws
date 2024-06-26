import { useCallback, useEffect, useState } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import Alert from "./components/Alert";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCircleUser } from '@fortawesome/free-solid-svg-icons';

function App() {

  const [jwtToken, setJwtToken] = useState("");
  const [userName, setUserName] = useState("");
  //const [roles, setRoles] = useState([]);
  const [isAdmin, SetIsAdmin] = useState(false);
  const [isUser, SetIsUser] = useState(false);
  const [alertMessage, setAlertMessage] = useState("");
  const [alertClassName, setAlertClassName] = useState("d-none");

  const [tickInterval, setTickInterval] = useState();

  const navigate = useNavigate();

  const logOut = () => {
    const requestOptions = {
      method: "GET",
      credentials: "include"
    }

    fetch("/logout", requestOptions)
      .catch( error => {
        console.log("error logging out", error)
      })
      .finally( () => {
        setJwtToken("")
        toggleRefresh(false)
      })
      setUserName("")
      //setRoles([])
      SetIsAdmin(false)
      SetIsUser(false)
      navigate("/login")
  }

  const toggleRefresh = useCallback((status) => {
    console.log("clicked")
    if (status) {
      console.log("turning on ticking")
      let i = setInterval( () => {
        console.log("this will run every second")
        const requestOptions = {
          method: "GET",
          credentials: 'same-origin'
        }
        fetch("/refresh", requestOptions)
          .then((response) => {
            return response.json()
          })
          .then((data) => {
            if (data.access_token) {
              setJwtToken(data.access_token)
            }
          })
          .catch( error => {
            console.log("user is not logged in")
          })
      }, 600000)
      setTickInterval(i)
      console.log("setting tick interval to ", i)
    } else {
      console.log("turning off ticking")
      console.log("turning off tickInterval", tickInterval)
      setTickInterval(null)
      clearInterval(tickInterval)
    }
  }, [tickInterval])

  useEffect( () => {
    if (jwtToken === "") {
      const requestOptions = {
        method: "GET",
        credentials: 'same-origin'
      }
      fetch("/refresh", requestOptions)
        .then((response) => {
          return response.json()
        })
        .then((data) => {
          if (data.access_token) {
            setJwtToken(data.access_token)
            toggleRefresh(true)
          }
        })
        .catch( error => {
          console.log("user is not logged in", error)
        })
    }
  }, [jwtToken, toggleRefresh])

  useEffect( () => {
    if (jwtToken !== "") {
      console.log(jwtToken)
      // Decode Token to get Username
      const base64Url = jwtToken.split('.')[1];
      const base64 = base64Url.replace('-', '+').replace('_', '/');
      const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));
      const { name, roles } = JSON.parse(jsonPayload);
      setUserName(name)
      //setRoles(roles)
      // Check if roles includes admin
      if (roles.includes("admin")) {
        SetIsAdmin(true)
      } else {
        SetIsAdmin(false)
      }
      // Check if roles includes user
      if (roles.includes("user")) {
        SetIsUser(true)
      }
    }
  }, [jwtToken])

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h1 className="mt-3"> Go Watch a Movie!</h1>
        </div>
        <div className="col d-flex flex-row-reverse align-items-center">
          {
            jwtToken ===""
            ? (
              <>
              <Link to="/login" ><span className="p-2 badge bg-success">Login</span></Link>
              <spacer></spacer>
              <Link to="/signup" ><span className="p-2 badge bg-info me-3">SignUp</span></Link>
              </>
            )
            : (
            <>
            <Link onClick={logOut} ><span className="p-2 badge bg-danger">Logout</span></Link>
            <p className="px-2 mt-3 text-success">{userName}</p>
            <FontAwesomeIcon icon={faCircleUser} size="lg" style={{color: "#63E6BE", }} ></FontAwesomeIcon>
            </>
          )
          }
        </div>
        <hr className="mb-3"></hr>
      </div>

      <div className="row">
        <div className="col-md-2">
          <nav>
            <div className="list-group">
              <Link to="/" className="list-group-item list-group-item-action">Home</Link>
              <Link to="/movies" className="list-group-item list-group-item-action">Movies</Link>
              <Link to="/genres" className="list-group-item list-group-item-action">Genres</Link>
              {
                jwtToken !== "" && isAdmin &&
                // Admin Links
                <>
                  <Link to="/admin/movie/0" className="list-group-item list-group-item-action">Add Movie</Link>
                  <Link to="/manage-catalogue" className="list-group-item list-group-item-action">Manage Catalogue</Link>
                  <Link to="/graphql" className="list-group-item list-group-item-action">GraphQL</Link>
                </>
              }
              
            </div>
          </nav>
        </div>
        <div className="col-md-10">
          <Alert
          message={alertMessage}
          className={alertClassName}
          />
          <Outlet context={{
            jwtToken,
            userName,
            isAdmin,
            isUser,
            setJwtToken,
            setAlertClassName,
            setAlertMessage,
            toggleRefresh
          }}></Outlet>
        </div>
      </div>
    </div>
  );
}

export default App;
