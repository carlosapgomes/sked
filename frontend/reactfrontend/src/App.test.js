import React from "react";
import { render } from "@testing-library/react";
import App from "./App";

it("renders without crashing", () => {
  const div = document.createElement("div");
  render(<App />, div);
});
//test('renders learn react link', () => {
//const { getByText } = render(<App />);
//const linkElement = getByText(/learn react/i);
//expect(linkElement).toBeInTheDocument();
//});
