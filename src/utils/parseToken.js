// src/utils/tokenParser.js
import jwt_decode from "jwt-decode";

export function parseToken(token) {
  try {
    const decoded = jwt_decode(token);
    return decoded;
  } catch (error) {
    console.error("TokenΩ‚Œˆ ß∞‹", error);
    return null;
  }
}