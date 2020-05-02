import {
  SET_ERRORS,
  SET_UNAUTHENTICATED,
  SET_AUTHENTICATED,

} from '../types';
import axios from 'axios';
import qs from 'qs';
import { hist } from '../../App'

export const loginUser = (userData) => (dispatch) => {
  const user = qs.stringify({ 'username': userData.username, 'password': userData.password });
  let myurl = 'http://localhost:8002/auth/login'
  axios({
    method: 'post',
    url: myurl,
    data: user,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
  }).then((res) => {
    if (res.data.error != null) {
      console.log('there are some error')
      alert(res.data.error)
      hist.push("/login")
      dispatch({
        type: SET_ERRORS,
      })
      return
    }

    dispatch({
      type: SET_AUTHENTICATED,
      payload: userData
    })
    hist.push("/admin")
    localStorage.setItem("isAuth", true)
  }).catch((err) => {
    hist.push("/login")
    dispatch({
      type: SET_ERRORS,
    })
  });
};

export const logoutUser = () => (dispatch) => {
  localStorage.setItem('isAuth', 'false')
  dispatch({ type: SET_UNAUTHENTICATED });
  hist.push("/login")
};
