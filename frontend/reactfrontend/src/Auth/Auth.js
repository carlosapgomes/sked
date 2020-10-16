import React, { useState } from "react";

const Auth = (props) => {
  const [username, setUserName] = useState("");
  const [pword, setPword] = useState("");
  const loginHandler = () => {
    let ajax = new XMLHttpRequest();
    ajax.open("POST", "https://myapp.local", true);
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.Stringfy({ username: username, password: pword }));
    ajax.onreadystatechange = () => {
      if (ajax.readyState == 4 && ajax.status == 200) {
        var data = ajax.responseText;
        console.log(data);
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
          placeholder="Username"
          onChange={(e) => setUserName(e.target.value)}
        />
        <input type="password" onChange={(e) => setPword(e.target.value)} />
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
