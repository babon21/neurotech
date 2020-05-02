import {
  SET_ERRORS,
  SET_UNAUTHENTICATED,
  SET_AUTHENTICATED,

} from '../types';
import axios from 'axios';
import qs from 'qs';
import {hist} from '../../App'

export const loginUser = (userData) => (dispatch) => {
  const user = qs.stringify({ 'username': userData.username, 'password': userData.password });
  alert('auth start')
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
      dispatch( {
        type: SET_ERRORS,
      })
    }

    dispatch( {
      type: SET_AUTHENTICATED,
      payload: userData
    })
    alert("auth success then /admin push")
    hist.push("/admin")
  }).catch((err) => {
    alert("auth catch fail")
    hist.push("/login")
    dispatch( {
      type: SET_ERRORS,
    })
  });
};

export const logoutUser = () => (dispatch) => {
  localStorage.removeItem('FBIdToken');
  delete axios.defaults.headers.common['Authorization'];
  dispatch({ type: SET_UNAUTHENTICATED });
};
