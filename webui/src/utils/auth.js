import Cookies from 'js-cookie'

const TokenKey = 'edge_admin_token'
const RefreshTokenKey = 'edge_admin_refresh_token'

export function getToken() {
  return Cookies.get(TokenKey)
}
export function getRefreshToken() {
  return Cookies.get(RefreshTokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token)
}

export function setRefreshToken(token) {
  return Cookies.set(RefreshTokenKey, token)
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function removeRefreshToken() {
  return Cookies.remove(RefreshTokenKey)
}