import React from "react";
import { render } from "@testing-library/react";
import ScheduleList from "./ScheduleList";
import i18n from "../i18nForTests";
import { I18nextProvider } from "react-i18next";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(
    <I18nextProvider i18n={i18n}>
      <ScheduleList />
    </I18nextProvider>,
    div
  );
});
