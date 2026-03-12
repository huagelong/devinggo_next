import React, { useState } from 'react';
import { LoginForm, ProFormText, ProFormCheckbox } from '@ant-design/pro-components';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { message } from 'antd';
import { useNavigate } from '@tanstack/react-router';
import { useAuthStore } from '../../stores/authStore';
import { login, getInfo } from '../../services/auth';
import { useTranslation } from 'react-i18next';

const LoginPage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const setToken = useAuthStore((state) => state.setToken);
  const setUserInfo = useAuthStore((state) => state.setUserInfo);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values: any) => {
    try {
      setLoading(true);
      // 调用登录接口
      const res = await login({
        username: values.username,
        password: values.password,
      });
      
      // 假设后端返回的数据包含 token
      if (res.token) {
        setToken(res.token);
        
        // 可选：登录成功后立刻拉取个人信息
        const userInfoRes = await getInfo();
        setUserInfo(userInfoRes);

        message.success('登录成功');
        navigate({ to: '/' });
      }
    } catch (error) {
      console.error(error);
      // 错误信息在 request 拦截器中已经抛出，无需在此重复
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ backgroundColor: 'white', height: '100vh', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
      <div style={{ width: '400px' }}>
        <LoginForm
          title="DevingGo-Light Admin"
          subTitle="高性能商业化管理后台"
          onFinish={handleSubmit}
          submitter={{
            searchConfig: {
              submitText: t('system.login'),
            },
            submitButtonProps: {
              loading: loading,
              size: 'large',
              style: { width: '100%' },
            },
          }}
        >
          <ProFormText
            name="username"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined className={'prefixIcon'} />,
            }}
            placeholder={'用户名'}
            rules={[{ required: true, message: '请输入用户名!' }]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined className={'prefixIcon'} />,
            }}
            placeholder={'密码'}
            rules={[{ required: true, message: '请输入密码!' }]}
          />
          <div style={{ marginBlockEnd: 24 }}>
            <ProFormCheckbox noStyle name="autoLogin">
              自动登录
            </ProFormCheckbox>
            <a style={{ float: 'right' }}>忘记密码 ?</a>
          </div>
        </LoginForm>
      </div>
    </div>
  );
};

export default LoginPage;
