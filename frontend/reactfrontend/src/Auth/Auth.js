import React, { useState } from "react";

const Auth = (props) => {
  const [email, setEmail] = useState("");
  const [pword, setPword] = useState("");
  const updateUser = (u) => {
    props.updateUser(u);
  };
  const updateLogin = (s) => {
    props.updateLogin(s);
  };

  const loginHandler = () => {
    let ajax = new XMLHttpRequest();
    ajax.open("POST", "https://dev.local/api/users/login", true);
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify({ email: email, password: pword }));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        var data = JSON.parse(ajax.responseText);
        if (typeof data.name != "undefined") {
          console.log("we got valid data");
          console.log(data);
          var user = {
            name: data.name,
            email: data.email,
          };
          updateUser(user);
          updateLogin(true);
        }
      }
    };
  };
  const logoutHandler = () => {};
  const cancelHandler = () => {};

  return (
    <div>
      <div id="login" hidden={props.loggedIn}>
        <input
          type="text"
          placeholder="Email"
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Senha"
          onChange={(e) => setPword(e.target.value)}
        />
        <button onClick={loginHandler}>Ok</button>
        <button onClick={cancelHandler}>Cancelar</button>
      </div>
      <div id="logout" hidden={!props.loggedIn}>
        <button onClick={logoutHandler}>Logout</button>
      </div>
    </div>
  );
};

export default Auth;
