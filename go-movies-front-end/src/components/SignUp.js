import Input from "./form/Input";
import Checkbox from "./form/Checkbox";
import {  useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Swal from 'sweetalert2'
const SignUp = () => {
    const navigate = useNavigate();
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    //constant for role with check status
    const [role, setRole] = useState({
        user: false,
        admin: false
    });
    // payloadrole constant
    const [roles, setRoles] = useState([])

    //const [ error, setError ] = useState(null);
    const [ errors, setErrors ] = useState([]);

    const hasError = (key) => {
        return errors.indexOf(key) !== -1;
    }

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

    // useEffect to set roles
    useEffect(() => {
        //console.log("useEffect called")
        let tempRoles = []
        //console.log("role state:", role)
        //console.log("roles state:", roles)
        // iterate over role state
        for (let key in role) {
            // if checked is true
            if (role[key]) {
                // push key to tempRoles as role tag

                tempRoles.push({"role": key})
            }
        }
        // set roles state with tempRoles
        setRoles(tempRoles)
        //console.log("roles:", roles)
    }, [role])


    const handleCheck = (event, position) => {
        //console.log("handleCheck called")
        //console.log("value in handleCheck:", event.target.value)
        //console.log("checked is", event.target.checked)
        //console.log("position is", position)
        // clone role state
        let tempRole = {...role}
        // set checked attribute in role state
        tempRole[event.target.value] = event.target.checked
        // set role state with new role
        setRole(tempRole)

    }


    const handleSubmit = (event) => {
        event.preventDefault();
        let errors = []
        let required = [
            { field: firstName, name: "firstname" },
            { field: lastName, name: "lastname" },
            { field: email, name: "email" },
            { field: password, name: "password" },
        ]

        required.forEach( function (obj) {
            if ( obj.field === ""){
                errors.push(obj.name)
            }
        })

        if (roles.length === 0) {
            alert("You must choose at least one role")
            errors.push("roles")
        }

        setErrors(errors);
        if ( errors.length > 0 ) {
            return false;
        }
        
        let payload = {
            first_name: firstName,
            last_name: lastName,
            email: email,
            password: password,
            roles: roles


        }

        console.log('payload',payload)

        const requestOptions = {
            method: "POST",
            headers: {
                "Content-Type":"application/json"
            },
            //credentials: 'include',
            credentials: 'same-origin',
            body: JSON.stringify(payload)
        }

        fetch("/signup", requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast.fire({
                        icon: "error",
                        title: "Failed to signup the user."
                    });
                    
                } else {
                    Toast.fire({
                        icon: "success",
                        title: "User onboarded successfully."
                    });
                    navigate("/")
                }
            })
            .catch(error => {
                console.log(error)
                Toast.fire({
                    icon: "error",
                    title: "Failed to signup the user."
                });
            })


    }

    return (
        <div className="col-md-6 offset-md-3">
            <h2>SignUp</h2>
            <hr></hr>

            <form onSubmit={handleSubmit}>
                <Input
                    title="First Name"
                    type="text"
                    className="form-control"
                    name="firstname"
                    autoComplete="firstname-new"
                    onChange={(event) => setFirstName(event.target.value)}
                    errorDiv={ hasError("firstname") ? "text-danger": "d-none" }
                    errorMsg={"Please enter your first name"}
                />
                <Input
                    title="Last Name"
                    type="text"
                    className="form-control"
                    name="lastname"
                    autoComplete="lastname-new"
                    onChange={(event) => setLastName(event.target.value)}
                    errorDiv={ hasError("lastname") ? "text-danger": "d-none" }
                    errorMsg={"Please enter your last name"}
                />
                <Input
                    title="Email Address"
                    type="email"
                    className="form-control"
                    name="email"
                    autoComplete="email-new"
                    onChange={(event) => setEmail(event.target.value)}
                    errorDiv={ hasError("email") ? "text-danger": "d-none" }
                    errorMsg={"Please enter your email"}
                />
                <Input
                    title="Password"
                    type="password"
                    className="form-control"
                    name="password"
                    autoComplete="password-new"
                    onChange={(event) => setPassword(event.target.value)}
                    errorDiv={ hasError("password") ? "text-danger": "d-none" }
                    errorMsg={"Please enter your password"}
                />
                <p>Roles</p>
                <Checkbox
                    title="user"
                    name={"userrole"}
                    key={1}
                    id={"user"}
                    onChange={(event) => {
                        handleCheck(event, 1)
                    }}
                    value={"user"}
                    checked={role.user}
                />
                <Checkbox
                    title={"admin"}
                    name={"adminrole"}
                    key={2}
                    id={"admin"}
                    onChange={(event) => {
                        handleCheck(event, 2)
                    }}
                    value={"admin"}
                    checked={role.admin}
                />
                <hr/>
                <input type="submit" className="btn btn-primary" value="Signup"></input>
            </form>
        </div>
    )
}

export default SignUp;