import React from "react";
import { render } from "@testing-library/react";
import Appointments from "./Appointments";
import { I18nextProvider } from "react-i18next";
import i18n from "../i18nForTests";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(
    <I18nextProvider i18n={i18n}>
      <Appointments doctors={["doc1", "doc2"]} />
    </I18nextProvider>,
    div
  );
});
