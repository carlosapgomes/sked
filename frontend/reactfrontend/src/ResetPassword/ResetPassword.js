import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import "./ResetPassword.css";

const ResetPassword = (props) => {
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const resetPassword = (e) => {
    e.preventDefault();
    let ajax = new XMLHttpRequest();
    ajax.open("POST", "/api/users/resetPassword", true);
    ajax.setRequestHeader("Content-type", "application/json");
    ajax.send(JSON.stringify({ email: email }));
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
      <section>
        <form>
          <p>{t("ResetPasswordMsg")}</p>
          <input
            type="text"
            placeholder="Email"
            onChange={(e) => setEmail(e.target.value)}
          />
          <button
            onClick={(e) => {
              resetPassword(e);
            }}
          >
            {t("SendLink")}
          </button>
        </form>
      </section>
    </div>
  );
};

export default ResetPassword;
