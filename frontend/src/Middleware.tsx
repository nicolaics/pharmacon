import { NavigateFunction, useNavigate } from "react-router-dom";
import { BACKEND_BASE_URL } from "./App";

export const AuthMiddleware = (navigate: NavigateFunction, admin: boolean) => {
  const token = sessionStorage.getItem("token");

  if (!token) {
    // No token found, redirect to login
    navigate("/");
  } else {
    const url = `http://${BACKEND_BASE_URL}/user/validate`;

    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + token,
      },
      body: JSON.stringify({
        needAdmin: admin,
      }),
    })
      .then((response) =>
        response.json().then((data) => {
          if (!response.ok) {
            throw new Error("Invalid credentials or network issue");
          }

          console.log("validated!");
        })
      )
      .catch((error) => {
        console.error("Error loading user data:", error);
        navigate("/home");
      });
  }
};
