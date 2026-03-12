import request from '../utils/request';

// API: /login POST
export const login = (data: any) => {
  return request.post('/login', data);
};

// API: /logout POST
export const logout = () => {
  return request.post('/logout');
};

// API: /refresh POST
export const refreshToken = () => {
  return request.post('/refresh');
};

// API: /getInfo GET (对应 user.go 获取自身信息)
export const getInfo = () => {
  return request.get('/getInfo');
};
