// author: j2.nete@gmail.com

using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Security.Cryptography;
using System.IO;

namespace tryExcel01
{
    class Program
    {
        static void Main(string[] args)
        {
            //TryECMA376();

            //TryRC4();

            //byte[] pwdVerifier = new byte[] { 0x46, 0x3F, 0xAC, 0xD7, 0x61, 0x05, 0x1B, 0xC1, 0x80, 0x09, 0x9A, 0xBD, 0x2F, 0x21, 0x03, 0x69 };
            //byte[] pwdVerifierHash = new byte[] { 0xCD, 0xEB, 0x07, 0x95, 0x85, 0xDB, 0x21, 0x60, 0xA2, 0xAE, 0x32, 0xDA, 0x08, 0x12, 0x8A, 0x79 };
            //byte[] key = new byte[] { 0xA6, 0xD9, 0x36, 0xB1, 0x11, 0x68, 0x75, 0xEF, 0x7F, 0x22, 0x6A, 0x22, 0x61, 0x86, 0x61, 0xC6 };

            byte[] pwdVerifier = new byte[] { 0x83, 0xfe, 0x18, 0xdf, 0x6c, 0xc3, 0xa5, 0x07, 0x23, 0x55, 0x31, 0x5f, 0xe5, 0xa2, 0xd2, 0x24 };
            byte[] pwdVerifierHash = new byte[] { 0x94, 0xec, 0xb4, 0x85, 0x15, 0x3d, 0xa2, 0x5b, 0xe2, 0x74, 0xe8, 0x6e, 0x53, 0x28, 0x35, 0xa8, 0xdd, 0xdf, 0x7e, 0x83 };
            byte[] key = new byte[] { 0xc4, 0xd9, 0x9c, 0x5c, 0x42, 0x3d, 0x02, 0x28, 0xc8, 0x2c, 0x46, 0x78, 0x43, 0xdd, 0x7f, 0xad };

            ARC4 arc4 = new ARC4();
            arc4.SetKey(key);

            byte[] pwdVerifierDec = arc4.En_De_Crypt(pwdVerifier);

            //arc4 = new ARC4();
            //arc4.SetKey(key);

            byte[] pwdVerifierHashDec = arc4.En_De_Crypt(pwdVerifierHash);

            arc4.consoleByte(pwdVerifierDec, "pwdVerifierDec");
            arc4.consoleByte(pwdVerifierHashDec, "pwdVerifierHashDec");

            Console.Read();
        }

        private static void TryRC4()
        {
            byte[] pwd = Encoding.Unicode.GetBytes("1");

            byte[] salt = new byte[] { 0xFF, 0xEF, 0xE3, 0x56, 0xD4, 0x5A, 0x2F, 0xBB, 0xC1, 0xD3, 0xFA, 0x60, 0xD3, 0x07, 0x2C, 0x54 };
            byte[] pwdVerifier = new byte[] { 0x9B, 0xFD, 0x09, 0x42, 0xAF, 0xFD, 0x6E, 0xB9, 0xE9, 0x1D, 0x40, 0x8D, 0x85, 0xFD, 0x85, 0x1D };
            byte[] pwdVerifierHash = new byte[] { 0xBF, 0x86, 0x91, 0x76, 0x84, 0x5D, 0xC5, 0x94, 0x18, 0xBA, 0xE4, 0xA1, 0x14, 0xE4, 0xC1, 0x4A };

            //byte[] salt = new byte[] { 0x13, 0xD3, 0x3B, 0x20, 0x89, 0xD5, 0x3F, 0xA9, 0x00, 0x68, 0x88, 0x0E, 0xE8, 0x04, 0xC6, 0x70 };
            //byte[] pwdVerifier = new byte[] { 0xF9, 0xE1, 0x1E, 0x5F, 0x95, 0xED, 0x4E, 0x70, 0xA0, 0xF6, 0xC2, 0x86, 0x6B, 0x6A, 0x64, 0x66 };
            //byte[] pwdVerifierHash = new byte[] { 0x88, 0xEB, 0xFE, 0xCC, 0xEA, 0x1D, 0x66, 0xF7, 0x50, 0xA8, 0xC2, 0xEA, 0xBC, 0x0B, 0xD6, 0x3D };

            MD5CryptoServiceProvider md5 = new MD5CryptoServiceProvider();
            byte[] pwdHash = md5.ComputeHash(pwd);

            byte[] tmp = new byte[21];
            Array.Copy(pwdHash, 0, tmp, 0, 5);
            Array.Copy(salt, 0, tmp, 5, salt.Length);

            byte[] tmpH = new byte[336];
            for (int i = 0; i < 16; i++)
            {
                tmp.CopyTo(tmpH, tmp.Length * i);
            }

            tmpH = md5.ComputeHash(tmpH);

            byte[] key = new byte[9];
            Array.Copy(tmpH, key, 5);

            Array.Copy(BitConverter.GetBytes(0), 0, key, 5, 4);

            key = md5.ComputeHash(key);

            ARC4 arc4 = new ARC4();
            arc4.SetKey(key);

            byte[] pwdVerifierDec = arc4.En_De_Crypt(pwdVerifier);
            byte[] pwdVerifierHashDec = arc4.En_De_Crypt(pwdVerifierHash);

            pwdVerifierDec = md5.ComputeHash(pwdVerifierDec);

            for (int i = 0; i < pwdVerifierDec.Length; i++)
            {
                if (pwdVerifierDec[i] != pwdVerifierHashDec[i])
                {
                    Console.WriteLine("Oh, NO!!!!!!!!!!!!!!!");
                    return;
                }
            }
            Console.WriteLine("Yahooooooooooooo~~~~~");
        }

        private static void TryECMA376()
        {
            byte[] salt = new byte[] { 0x30, 0x3A, 0x39, 0x9E, 0x02, 0xB6, 0x66, 0x71, 0xA9, 0xEA, 0xDB, 0x17, 0x0C, 0xA6, 0xDF, 0x24 };
            byte[] pwd = Encoding.Unicode.GetBytes("123");
            byte[] tmp = H(salt, pwd);
            for (int i = 0; i < 50000; i++)
            {
                byte[] tmpI = BitConverter.GetBytes(i);
                tmp = H(tmpI, tmp);
            }
            tmp = H(tmp, BitConverter.GetBytes(0));

            SHA1Managed sha = new SHA1Managed();
            byte[] XorBuff = new byte[64];
            for (int i = 0; i < XorBuff.Length; i++)
            {
                XorBuff[i] = 0x36;
            }

            for (int i = 0; i < tmp.Length; i++)
            {
                XorBuff[i] ^= tmp[i];
            }

            tmp = sha.ComputeHash(XorBuff);

            byte[] key = new byte[16];
            Array.Copy(tmp, key, 16);

            byte[] pwdVerifier = new byte[] { 0x0A, 0x3A, 0xD5, 0x44, 0x17, 0xAB, 0x9B, 0x26, 0xA7, 0xFC, 0x65, 0xFE, 0x2F, 0x77, 0xDD, 0xA0 };
            byte[] pwdVerifierHash = new byte[] { 0x29, 0x0E, 0x48, 0xC6, 0x26, 0xA8, 0x72, 0x3B, 0xA2, 0xCF, 0x3F, 0x6A, 0x9B, 0x3A, 0x94, 0xB8, 0x59, 0x7A, 0x8E, 0xCE, 0xB5, 0x55, 0xE9, 0x15, 0x06, 0x89, 0xCD, 0x9C, 0xDA, 0x6C, 0x39, 0x31 };

            AesManaged aes = new AesManaged();
            aes.Mode = CipherMode.ECB;
            aes.KeySize = 128;
            aes.Key = key;
            aes.Padding = PaddingMode.None;
            ICryptoTransform ict = aes.CreateDecryptor();
            ICryptoTransform ict2 = aes.CreateDecryptor();
            byte[] pwdVerifierDec;
            byte[] pwdVerifierHashDec;
            pwdVerifierDec = ict.TransformFinalBlock(pwdVerifier, 0, pwdVerifier.Length);
            pwdVerifierHashDec = ict.TransformFinalBlock(pwdVerifierHash, 0, pwdVerifierHash.Length);
            //pwdVerifierDec = new byte[16];
            //pwdVerifierHashDec = new byte[32];
            //ict.TransformBlock(pwdVerifier, 0, pwdVerifier.Length, pwdVerifierDec, 0);
            //ict2.TransformBlock(pwdVerifierHash, 0, 32, pwdVerifierHashDec, 0);
            //ict2.TransformBlock(pwdVerifierHash, 16, 16, pwdVerifierHashDec, 16);

            pwdVerifierDec = sha.ComputeHash(pwdVerifierDec);

            for (int i = 0; i < pwdVerifierDec.Length; i++)
            {
                if (pwdVerifierDec[i] != pwdVerifierHashDec[i])
                {
                    Console.WriteLine("Oh, NO!!!!!!!!!!!!!!!");
                    return;
                }
            }
            Console.WriteLine("Yahooooooooooooo~~~~~");

            FileStream fs = File.OpenRead("C:\\Users\\Administrator\\Desktop\\加密\\EncryptedPackage");
            byte[] encPak = new byte[fs.Length];
            fs.Read(encPak, 0, encPak.Length);
            byte[] decPak = new byte[fs.Length];
            fs.Close();
            decPak = ict2.TransformFinalBlock(encPak, 8, encPak.Length - 8);
            FileStream fso = File.Create("C:\\Users\\Administrator\\Desktop\\加密\\yeah.zip");
            fso.Write(decPak, 0, decPak.Length);
            fso.Close();
        }

        private static byte[] H(byte[] b1, byte[] b2)
        {
            SHA1Managed sha = new SHA1Managed();
            byte[] t = new byte[b1.Length + b2.Length];
            b1.CopyTo(t, 0);
            b2.CopyTo(t, b1.Length);
            return sha.ComputeHash(t);
        }
    }

    class ARC4
    {
        public void consoleByte(byte[] bb, string str)
        {
            Console.WriteLine(str);
            foreach (byte b in bb)
            {
                Console.Write("{0:X} ", b);
            };
            Console.WriteLine();

        }

        private int m_x, m_y;
        private byte[] m_State = new byte[256];

        public void SetKey(byte[] key)
        {
            m_x = 0;
            m_y = 0;
            for (int i = 0; i < 256; i++)
                m_State[i] = (byte)i;
            int k_len = key.Length;
            int index1 = 0;
            int index2 = 0;
            for (int i = 0; i < 256; i++)
            {
                unchecked
                {
                    index2 = (byte)(key[index1] + m_State[i] + index2);
                }
                byte t = m_State[i];
                m_State[i] = m_State[index2];
                m_State[index2] = t;
                index1 = (byte)((index1 + 1) % k_len);
            }
        }

        public byte[] En_De_Crypt(byte[] input)
        {
            if (input == null || input.Length <= 0) return null;

            int x = m_x;
            int y = m_y;
            byte[] output = new byte[input.Length];

            for (int i = 0; i < input.Length; i++)
            {
                byte c = input[i];
                int x1 = x = (byte)(x + 1);
                int y1 = y = (byte)(m_State[x] + y);
                byte tx = m_State[x1];
                byte ty = m_State[y1];
                m_State[x1] = ty;
                m_State[y1] = tx;
                output[i] = (byte)(c ^ m_State[(byte)(tx + ty)]);
            }

            m_x = x;
            m_y = y;

            return output;
        }
    }
}
