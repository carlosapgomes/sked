import React from "react";
import { render } from "@testing-library/react";
import App from "./App";
import { I18nextProvider } from "react-i18next";
import i18n from "./i18nForTests";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(
    <I18nextProvider i18n={i18n}>
      <App />
    </I18nextProvider>,
    div
  );
});
//test('renders learn react link', () => {
//const { getByText } = render(<App />);
//const linkElement = getByText(/learn react/i);
//expect(linkElement).toBeInTheDocument();
//});
