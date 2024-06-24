import { useEffect, useState } from "react";
import Input from "./form/Input";
import { useNavigate, useOutletContext } from "react-router-dom";

const Login = () => {

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const { jwtToken, setJwtToken, setAlertClassName, setAlertMessage, toggleRefresh} = useOutletContext();
    const navigate = useNavigate();


    useEffect(() => {
        setAlertClassName("d-none");
        setAlertMessage("");
        if (jwtToken) {
            navigate("/")
        }
    }, [jwtToken, navigate])


    const handleSubmit = (event) => {
        event.preventDefault();
        // console.log("email/password",email, password)
        // if ( email === "admin@example.com") {
        //     setJwtToken(email)
        //     setAlertClassName("d-none")
        //     setAlertMessage("")
        //     navigate("/")
        // } else {
        //     setAlertClassName("alert-danger")
        //     setAlertMessage("Invalid Credential")
        // }
        let payload = {
            email: email,
            password: password
        }

        const requestOptions = {
            method: "POST",
            headers: {
                "Content-Type":"application/json"
            },
            //credentials: 'include',
            credentials: 'same-origin',
            body: JSON.stringify(payload)
        }

        fetch("/authenticate", requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    setAlertClassName("alert-danger");
                    setAlertMessage(data.message);
                    setJwtToken("")
                } else {
                    setJwtToken(data.access_token);
                    setAlertClassName("d-none");
                    setAlertMessage("");
                    toggleRefresh(true)
                    navigate("/")
                }
            })
            .catch(error => {
                setAlertClassName("alert-danger")
                setAlertMessage(error)
            })


    }
    return(
        <div className="col-md-6 offset-md-3">
            <h2>Login</h2>
            <hr></hr>

            <form onSubmit={handleSubmit}>
                <Input
                    title="Email Address"
                    type="email"
                    className="form-control"
                    name="email"
                    autoComplete="email-new"
                    onChange={(event) => setEmail(event.target.value)}
                />
                <Input
                    title="Password"
                    type="password"
                    className="form-control"
                    name="password"
                    autoComplete="password-new"
                    onChange={(event) => setPassword(event.target.value)}
                />
                <hr/>
                <input type="submit" className="btn btn-primary" value="Login"></input>
            </form>
        </div>
    )
}

export default Login;