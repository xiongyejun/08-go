using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace OfficeRC4
{
    class Program
    {
        static void Main(string[] args)
        {
            System.IO.FileStream stream = new System.IO.FileInfo("C:/Users/Administrator/Desktop/rc4EncryptionInfo").OpenRead();
            byte[] buffer = new byte[stream.Length + 1];
            stream.Read(buffer, 0, Convert.ToInt32(stream.Length));

            CDecrypt decrypt = new CDecrypt();
            Console.WriteLine(decrypt.DecryptInternal("1", buffer));

            Console.Read();
        }
    }
}
