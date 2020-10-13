import { useState } from "react";

const Auth = (props) => {
  const [username, setUserName] = useState("");
  const [pword, setPword] = useState("");
  const loginHandler = () => {};
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
        <button onClick={this.loginHandler}>Ok</button>
        <button onClick={this.cancelHandler}>Cancelar</button>
      </div>
      <div id="logout" hidden={!props.loggedIn}>
        <button onClick={this.logoutHandler}>Logout</button>
      </div>
    </div>
  );
};

export default Auth;
