import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import "./Auth.css";

const Auth = (props) => {
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [pword, setPword] = useState("");
  const loginHandler = (e) => {
    e.preventDefault();
    let ajax = new XMLHttpRequest();
    ajax.open("POST", "/api/users/login", true);
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify({ email: email, password: pword }));
    ajax.onreadystatechange = () => {
      if (ajax.readyState === 4 && ajax.status === 200) {
        var data = JSON.parse(ajax.responseText);
        if (typeof data.name != "undefined") {
          var user = {
            name: data.name,
            email: data.email,
            uid: data.id,
            phone: data.phone,
            roles: data.roles,
          };
          props.updateUser(user);
          props.updateLogin(data.active && data.emailWasValidated);
        }
      }
      if (ajax.readyState === 4 && ajax.status !== 200) {
        console.log(ajax.responseText);
        window.alert("Could not login");
      }
    };
  };

  return (
    <div>
      <div id="login" hidden={props.loggedIn}>
        <section>
          <form>
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
            <button
              onClick={(e) => {
                loginHandler(e);
              }}
            >
              Ok
            </button>
            <p>
              {t("ForgottenPw")}
              <Link className="Link" to="/ResetPassword">
                &nbsp;
                {t("ClickHere")}
              </Link>
            </p>
          </form>
        </section>
      </div>
    </div>
  );
};

export default Auth;
