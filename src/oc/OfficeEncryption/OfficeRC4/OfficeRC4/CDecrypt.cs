using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace OfficeRC4
{
    class CDecrypt
    {
        byte[] encryptedVerifier = null;
        byte[] encryptedVerifierHash = null;
        byte[] salt;
        UInt32 keySize;
        Int32 verifierHashSize;

        #region Derive key

        private byte[] DeriveKey(byte[] hashValue)
        {
            // And one more hash to derive the key
            byte[] derivedKey = new byte[64];

            // This is step 4a in 2.3.4.7 of MS_OFFCRYPT version 1.0
            // and is required even though the notes say it should be 
            // used only when the encryption algorithm key > hash length.
            for (int i = 0; i < derivedKey.Length; i++)
                derivedKey[i] = (byte)(i < hashValue.Length ? 0x36 ^ hashValue[i] : 0x36);

            consoleByte(derivedKey, "异或36");

            byte[] X1 = SHA1Hash(derivedKey);
            consoleByte(X1, "X1");
            Console.WriteLine("keySize={0}", keySize);
            if (verifierHashSize > keySize / 8)
                return X1;

            for (int i = 0; i < derivedKey.Length; i++)
                derivedKey[i] = (byte)(i < hashValue.Length ? 0x5C ^ hashValue[i] : 0x5C);

            byte[] X2 = SHA1Hash(derivedKey);

            byte[] X3 = new byte[X1.Length + X2.Length];

            Array.Copy(X1, 0, X3, 0, X1.Length);
            Array.Copy(X1, 0, X3, X1.Length, X2.Length);

            consoleByte(X3, "x3");

            return X3;
        }

        #endregion

        byte[] GeneratePasswordHashUsingSHA1(string password)
        {
            byte[] hashBuf = null;

            try
            {
                // H(0) = H(salt, password);
                hashBuf = SHA1Hash(salt, password);

                for (int i = 0; i < 50000; i++)
                {
                    // Generate each hash in turn
                    // H(n) = H(i, H(n-1))
                    hashBuf = SHA1Hash(i, hashBuf);
                }

                // Finally, append "block" (0) to H(n)
                hashBuf = SHA1Hash(hashBuf, 0);
                consoleByte(hashBuf, "hashBuf");

                // The algorithm in this 'DeriveKey' function is the bit that's not clear from the documentation
                byte[] key = DeriveKey(hashBuf);

                // Should handle the case of longer key lengths as shown in 2.3.4.9
                // Grab the key length bytes of the final hash as the encrypytion key
                byte[] final = new byte[keySize / 8];
                Array.Copy(key, final, final.Length);

                consoleByte(final, "final");

                return final;

            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }

            return null;
        }

        byte[] AESDecrypt(byte[] data, int index, int count, byte[] key)
        {
            byte[] decryptedBytes = null;

            //  Create uninitialized Rijndael encryption object.
            System.Security.Cryptography.RijndaelManaged symmetricKey = new System.Security.Cryptography.RijndaelManaged();

            // It is required that the encryption mode is Electronic Codebook (ECB) 
            // see MS-OFFCRYPTO v1.0 2.3.4.7 pp 39.
            symmetricKey.Mode = System.Security.Cryptography.CipherMode.ECB;
            symmetricKey.Padding = System.Security.Cryptography.PaddingMode.None;
            symmetricKey.KeySize = 128;
            // symmetricKey.IV = null; // new byte[16];
            // symmetricKey.Key = key;

            //  Generate decryptor from the existing key bytes and initialization 
            //  vector. Key size will be defined based on the number of the key 
            //  bytes.
            System.Security.Cryptography.ICryptoTransform decryptor;
            decryptor = symmetricKey.CreateDecryptor(key, null);

            //  Define memory stream which will be used to hold encrypted data.
            using (System.IO.MemoryStream memoryStream = new System.IO.MemoryStream(data, index, count))
            {
                //  Define memory stream which will be used to hold encrypted data.
                using (System.Security.Cryptography.CryptoStream cryptoStream
                        = new System.Security.Cryptography.CryptoStream(memoryStream, decryptor, System.Security.Cryptography.CryptoStreamMode.Read))
                {
                    //  Since at this point we don't know what the size of decrypted data
                    //  will be, allocate the buffer long enough to hold ciphertext;
                    //  plaintext is never longer than ciphertext.
                    decryptedBytes = new byte[data.Length];
                    int decryptedByteCount = cryptoStream.Read(decryptedBytes, 0, decryptedBytes.Length);

                    return decryptedBytes;
                }
            }
        }

        private byte[] SHA1Hash(byte[] salt, string password)
        {
            return SHA1Hash(HashPassword(salt, password));
        }

        private byte[] HashPassword(byte[] salt, string password)
        {
            // Use a unicode form of the password
            byte[] passwordBuf = System.Text.UnicodeEncoding.Unicode.GetBytes(password);
            byte[] inputBuf = new byte[salt.Length + passwordBuf.Length];
            Array.Copy(salt, inputBuf, salt.Length);
            Array.Copy(passwordBuf, 0, inputBuf, salt.Length, passwordBuf.Length);

            return inputBuf;
        }

        private byte[] SHA1Hash(int iterator, byte[] hashBuf)
        {
            // Create an input buffer for the hash.  This will be 4 bytes larger than 
            // the hash to accommodate the unsigned int iterator value.
            byte[] inputBuf = new byte[0x14 + 0x04];

            // Create a byte array of the integer and put at the front of the input buffer
            // 1.3.6 says that little-endian byte ordering is expected

            // Copy the iterator bytes into the hash input buffer
            Array.Copy(System.BitConverter.GetBytes(iterator), inputBuf, 4);

            // 'append' the previously generated hash to the input buffer
            Array.Copy(hashBuf, 0, inputBuf, 4, hashBuf.Length);

            return SHA1Hash(inputBuf);
        }

        byte[] SHA1Hash(byte[] b)
        {
            System.Security.Cryptography.SHA1 sha1 = System.Security.Cryptography.SHA1CryptoServiceProvider.Create();
            return sha1.ComputeHash(b);
        }

        private byte[] SHA1Hash(byte[] hashBuf, int block)
        {
            // Create an input buffer for the hash.  This will be 4 bytes larger than 
            // the hash to accommodate the unsigned int iterator value.
            byte[] inputBuf = new byte[0x14 + 0x04];

            Array.Copy(hashBuf, inputBuf, hashBuf.Length);
            Array.Copy(System.BitConverter.GetBytes(block), 0, inputBuf, hashBuf.Length, 4);

            return SHA1Hash(inputBuf);
        }

        private byte[] SHA1Hash(byte[] hashBuf, byte[] block0)
        {
            // Create an input buffer for the hash.  This will be 4 bytes larger than 
            // the hash to accommodate the unsigned int iterator value.
            byte[] inputBuf = new byte[hashBuf.Length + block0.Length];

            Array.Copy(hashBuf, inputBuf, hashBuf.Length);
            Array.Copy(block0, 0, inputBuf, hashBuf.Length, block0.Length);

            return SHA1Hash(inputBuf);
        }

        bool PasswordVerifier(byte[] key)
        {
            // Decrypt the encrypted verifier...
            consoleByte(encryptedVerifier, "encryptedVerifier");
            byte[] decryptedVerifier = AESDecrypt(encryptedVerifier,0, encryptedVerifier.Length, key);
            
            consoleByte(decryptedVerifier, "decryptedVerifier");

            // Truncate
            byte[] data = new byte[16];
            Array.Copy(decryptedVerifier, data, data.Length);
            decryptedVerifier = data;

            // ... and hash
            byte[] decryptedVerifierHash = AESDecrypt(encryptedVerifierHash, 0, encryptedVerifierHash.Length, key);

            // Hash the decrypted verifier (2.3.4.9)
            byte[] checkHash = SHA1Hash(decryptedVerifier);

            // Check the 
            for (int i = 0; i < checkHash.Length; i++)
            {
                if (decryptedVerifierHash[i] != checkHash[i])
                    return false;
            }
            return true;
        }

    void  consoleByte(byte[]bb, string str)
        {
            Console.WriteLine(str);
            foreach (byte b in bb)
            {
                Console.Write("{0:x} ", b);
            };
            Console.WriteLine();

        }

    public    bool DecryptInternal(string password, byte[] encryptionInfo)
        {
   
            #region Parse the encryption info data

            using (System.IO.MemoryStream ms = new System.IO.MemoryStream(encryptionInfo))
            {
                System.IO.BinaryReader reader = new System.IO.BinaryReader(ms);

                ushort versionMajor = reader.ReadUInt16();
                ushort versionMinor = reader.ReadUInt16();

                UInt32 encryptionFlags = reader.ReadUInt32();
                if (encryptionFlags == 16)
                    throw new Exception("An external cryptographic provider is not supported");

                // Encryption header
                uint headerLength = reader.ReadUInt32();
                int skipFlags = reader.ReadInt32(); headerLength -= 4;
                UInt32 sizeExtra = reader.ReadUInt32(); headerLength -= 4;
                UInt32 algId = reader.ReadUInt32(); headerLength -= 4;
                UInt32 algHashId = reader.ReadUInt32(); headerLength -= 4;
                keySize = reader.ReadUInt32(); headerLength -= 4;
                UInt32 providerType = reader.ReadUInt32(); headerLength -= 4;
                reader.ReadUInt32(); headerLength -= 4; // Reserved 1
                reader.ReadUInt32(); headerLength -= 4; // Reserved 2
                string CSPName = System.Text.UnicodeEncoding.Unicode.GetString(reader.ReadBytes((int)headerLength));

                // Encryption verifier
                Int32 saltSize = reader.ReadInt32();
                salt = reader.ReadBytes(saltSize);
                encryptedVerifier = reader.ReadBytes(0x10);
                verifierHashSize = reader.ReadInt32();
                encryptedVerifierHash = reader.ReadBytes(providerType == 1 ? 0x14 : 0x20);

                consoleByte(encryptedVerifier, "encryptedVerifier");
                consoleByte(encryptedVerifierHash, "encryptedVerifierHash");
            }
            #endregion

            #region Encryption key generation

            Console.WriteLine("Encryption key generation");
            byte[] encryptionKey = GeneratePasswordHashUsingSHA1(password);
            if (encryptionKey == null) return false;

            #endregion

            #region Password verification

            Console.WriteLine("Password verification");
            if (PasswordVerifier(encryptionKey))
            {
                Console.WriteLine("Password verification succeeded");
                return true;
            }
            else
            {
                Console.WriteLine("Password verification failed");
                return false;
            }

            #endregion
                          }

    }
}
