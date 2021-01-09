import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import "./ResetPassword.css";

const ResetPassword = () => {
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
        window.alert("Request sent");
      }
    };
    if (ajax.readyState === 4 && ajax.status !== 200) {
      console.log(ajax.responseText);
      window.alert("Could send request");
    }
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
