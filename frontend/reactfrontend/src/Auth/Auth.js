import React, { useState } from "react";

const Auth = (props) => {
  const [email, setEmail] = useState("");
  const [pword, setPword] = useState("");
  const loginHandler = () => {
    let ajax = new XMLHttpRequest();
    ajax.open("POST", "https://dev.local/api/users/login", true);
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify({ email: email, password: pword }));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        var data = JSON.parse(ajax.responseText);
        if (typeof data.name != "undefined") {
          var user = {
            name: data.name,
            email: data.email,
          };
          props.updateUser(user);
          props.updateLogin(data.loggedin);
        }
      }
    };
  };

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
      </div>
    </div>
  );
};

export default Auth;
