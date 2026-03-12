import CryptoJS from 'crypto-js';

// AES 加密配置 (优先从环境变量读取，默认与后端保持一致)
const AES_KEY = import.meta.env.VITE_APP_AES_KEY || '1234567890123456';

/**
 * AES-ECB 加密
 * @param plaintext 明文
 * @returns Base64 编码的密文
 */
export const aesEncrypt = (plaintext: string): string => {
  const keyHex = CryptoJS.enc.Utf8.parse(AES_KEY);
  const encrypted = CryptoJS.AES.encrypt(plaintext, keyHex, {
    mode: CryptoJS.mode.ECB,
    padding: CryptoJS.pad.Pkcs7,
  });
  return encrypted.toString();
};
