using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.IO;
using System.Runtime.Serialization.Formatters.Binary;
//using ProtoBuf;
using Protocol.Dto;
using UnityEngine;
using Protocol.Code;

public class EncodeTool
{
    /// <summary>
    /// 构造包 包头+包尾
    /// </summary>
    /// <param name="data"></param>
    /// <returns></returns>
    public static byte[] EncodePacket(byte[] data)
    {
        using (MemoryStream ms = new MemoryStream())
        {
            using (BinaryWriter bw = new BinaryWriter(ms))
            {
                //写入包头（数据的长度）
                //bw.Write(data.Length);
                ////Console.Write(data.Length);
                ////写入包尾（数据）
                //bw.Write(data);

                /*测试 2019-4-1*/
                Int32 length = data.Length;
                byte[] arr2 = BitConverter.GetBytes((Int32)(length));


                if (BitConverter.IsLittleEndian)  //如果为xiaoduan
                {
                    Array.Reverse(arr2);
                }

                //byte[] packet = new byte[ms.Length];
                //Buffer.BlockCopy(ms.GetBuffer(), 0, packet, 0, (int)ms.Length);

                /*测试 2019-4-1*/
                byte[] packet = new byte[arr2.Length + data.Length];
                arr2.CopyTo(packet, 0);
                data.CopyTo(packet, arr2.Length);

                return packet;
            }
        }
    }
    /// <summary>
    /// 解析包，从缓冲区里取出一个完整的包
    /// </summary>
    /// <param name="cache"></param>
    /// <returns></returns>
    public static byte[] DecodePacket(ref List<byte> cache)
    {
        Debug.Log("DecodePacket");
        if (cache.Count < 4)
        {
            return null;
        }
        using (MemoryStream ms = new MemoryStream(cache.ToArray()))
        {
            using (BinaryReader br = new BinaryReader(ms))
            {

                //4-8 @author wangfw
                //原来是实际的读入，改写后，这里只是偏移
                  int length2 = (int)br.ReadUInt32();
                //获取包头
                byte[] arr1 = new byte[4];
                arr1[0] = cache[0];
                arr1[1] = cache[1];
                arr1[2] = cache[2];
                arr1[3] = cache[3];

                if (BitConverter.IsLittleEndian)
                {
                    Array.Reverse(arr1);
                }
                int length = (int)BitConverter.ToUInt32(arr1, 0);           
                
                //@wangfw 有问题
                int remainLength = (int)(ms.Length - ms.Position);
                if (length > remainLength)
                {
                    return null;
                }
                byte[] data = cache.ToArray();

                //构造值体
                byte[] packaet = data.Skip(4).ToArray();

                //原来是实际的读入，改写后，这里只是偏移
                byte[] data2 = br.ReadBytes(length2);

                //更新数据缓存
                cache.Clear();
                cache.AddRange(br.ReadBytes(remainLength));

           
                return packaet;
            }
        }
    }

    /// <summary>
    /// 把NetMsg类转换成字节数组，发送出去
    /// </summary>
    /// <param name="msg"></param>
    /// <returns></returns>
    public static byte[] EncodeMsg(NetMsg msg)
    {
        using (MemoryStream ms = new MemoryStream())
        {
            using (BinaryWriter bw = new BinaryWriter(ms))
            {
                bw.Write(msg.opCode);
                bw.Write(msg.subCode);
                if (msg.value != null)
                {
                    bw.Write(EncodeObj(msg.value));
                }
                byte[] data = new byte[ms.Length];
                Buffer.BlockCopy(ms.GetBuffer(), 0, data, 0, (int)ms.Length);
                return data;
            }
        }

        /*王付伟 2019-4-1 测试 */
        //Int32 Opcode = msg.opCode;
        //Int32 Subcode = msg.subCode;
        //byte[] arr1 = BitConverter.GetBytes((Int32)(Opcode));
        //byte[] arr2 = BitConverter.GetBytes((Int32)Subcode);

        //byte[] packet = new byte[arr1.Length + arr2.Length];

    }
    /// <summary>
    /// 将字节数组转换成 NetMsg 网络消息类
    /// </summary>
    /// <param name="data"></param>
    /// <returns></returns>
    public static NetMsg DecodeMsg(byte[] data)
    {
        using (MemoryStream ms = new MemoryStream(data))
        {
            using (BinaryReader br = new BinaryReader(ms))
            {
                NetMsg msg = new NetMsg();


                //4-8 wangfw
                byte[] opcode = new byte[4];
                opcode[0] = data[0];
                opcode[1] = data[1];
                opcode[2] = data[2];
                opcode[3] = data[3];
                if (BitConverter.IsLittleEndian)
                {
                    Array.Reverse(opcode);
                }
                msg.opCode = (int)BitConverter.ToUInt32(opcode, 0);

                byte[] subcode = new byte[4];
                subcode[0] = data[4];
                subcode[1] = data[5];
                subcode[2] = data[6];
                subcode[3] = data[7];
                if (BitConverter.IsLittleEndian)
                {
                    Array.Reverse(subcode);
                }
                msg.subCode = (int)BitConverter.ToUInt32(subcode, 0);

                byte[] value = data.Skip(8).ToArray();

                // msg.opCode = br.ReadInt32();
                // msg.subCode = br.ReadInt32();
                //判断是否还有value的值
                //new判断是否还有value的值  --wangfw 4-9

                if(value != null && value.Length != 0)
                {
                    object obj = NewDecodeoObj(msg.opCode,msg.subCode,value);
                    msg.value = obj;
                }
                else
                {
                    msg.value = null;

                }

                //if (ms.Length - ms.Position > 0)
                //{
                //    //object obj = DecodeObj(br.ReadBytes((int)(ms.Length - ms.Position)));
                //    object obj = NewDecodeoObj(value);
                //    msg.value = obj;
                //}

                // Debug.Log(msg.opCode);
                //  Debug.Log(msg.subCode);

                return msg;
            }
        }
    }
    /// <summary>
    /// 序列化
    /// </summary>
    /// <param name="obj"></param>
    /// <returns></returns>
    private static byte[] EncodeObj(object obj)
    {
        using (MemoryStream ms = new MemoryStream())
        {
            BinaryFormatter bf = new BinaryFormatter();
            bf.Serialize(ms, obj);
            byte[] data = new byte[ms.Length];
            Buffer.BlockCopy(ms.GetBuffer(), 0, data, 0, (int)ms.Length);
            return data;
        }
    }
    /// <summary>
    /// 反序列化
    /// </summary>
    /// <param name="data"></param>
    /// <returns></returns>
    private static object DecodeObj(byte[] data)
    {
        using (MemoryStream ms = new MemoryStream(data))
        {
            BinaryFormatter bf = new BinaryFormatter();
            return bf.Deserialize(ms);
        }
    }

    /// <summary>
    /// 反序列化
    /// @author wangfw
    /// @time 2019-4-4
    /// </summary>
    /// <param name="data"></param>
    /// <returns></returns>
    public static object NewDecodeoObj(int opcode,int subcode,byte[] data)
    {
        //函数可以扩展参数，subcode和supcode来判断data的类型
      

        //value为userDto
        if ((OpCode.Account == opcode && AccountCode.GetUserInfo_SRES == subcode) ||
            (OpCode.Match == opcode && MatchCode.Enter_BRO == subcode))   //GetUserInfo_SRES UserDto
        {
            //最终格式为 len+value格式（id+ len + name + len + InconName + jinbi）len为uint32

            // List<byte> mysrc = new List<byte>(data);

            UserDto dto = NewUpackUserDto(data);
            return dto;
        }
        else if(OpCode.Match == opcode && MatchCode.Enter_SRES == subcode)
        {
            //格式：len+map+ len+list+ len+list
            byte[] arr1 = new byte[4];

            arr1[0] = data[0];
            arr1[1] = data[1];
            arr1[2] = data[2];
            arr1[3] = data[3];
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(arr1);
            }

            int lenmap = (int)BitConverter.ToUInt32(arr1, 0);
            byte[] dataMAP = data.Skip(4).Take(4 + lenmap).ToArray();
            data = data.Skip(4 + lenmap).ToArray();
            Dictionary<int, UserDto> userIdUserDtoDic = NewUpackuserIdUserDtoDic(dataMAP);

            byte[] arr2 = new byte[4];
            arr2[0] = data[0];
            arr2[1] = data[1];
            arr2[2] = data[2];
            arr2[3] = data[3];
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(arr2);
            }
            int lenList = (int)BitConverter.ToUInt32(arr2, 0);
            byte[] dataList = data.Skip(4).Take(4 + lenList).ToArray();
            data = data.Skip(4 + lenList).ToArray();
            List<int> readyUserIdList = NewUpackListInt(dataList, lenList / 4);

            byte[] arr3 = new byte[4];
            arr3[0] = data[0];
            arr3[1] = data[1];
            arr3[2] = data[2];
            arr3[3] = data[3];
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(arr3);
            }

            int lenList2 = (int)BitConverter.ToUInt32(arr2, 0);
            byte[] dataList2 = data.Skip(4).Take(4 + lenList2).ToArray();
            List<int> enterOrderUserIdList = NewUpackListInt(dataList2, lenList2 / 4);

            MatchRoomDto2 Room = new MatchRoomDto2();
            Room.readyUserIdList = enterOrderUserIdList;
            Room.userIdUserDtoDic = userIdUserDtoDic;
            Room.enterOrderUserIdList = enterOrderUserIdList;
            return Room;
        }

        //value为一个int
        object value = new object();
        if (BitConverter.IsLittleEndian)
        {
            Array.Reverse(data);
        }
        value = BitConverter.ToInt32(data, 0);

        return value;
    }

    private static UserDto NewUpackUserDto(byte[] data)
    {
        byte[] arr1 = new byte[4];

        arr1[0] = data[0];
        arr1[1] = data[1];
        arr1[2] = data[2];
        arr1[3] = data[3];
        if (BitConverter.IsLittleEndian)
        {
            Array.Reverse(arr1);
        }

        int id = (int)BitConverter.ToUInt32(arr1, 0);
        data = data.Skip(4).ToArray();

        arr1[0] = data[0];
        arr1[1] = data[1];
        arr1[2] = data[2];
        arr1[3] = data[3];
        if (BitConverter.IsLittleEndian)
        {
            Array.Reverse(arr1);
        }

        int lenN = (int)BitConverter.ToUInt32(arr1, 0);
        byte[] namebuf = new byte[lenN];
        namebuf = data.Skip(4).Take(lenN).ToArray();
        string name = Encoding.UTF8.GetString(namebuf);
        data = data.Skip(4 + lenN).ToArray();

        arr1[0] = data[0];
        arr1[1] = data[1];
        arr1[2] = data[2];
        arr1[3] = data[3];
        if (BitConverter.IsLittleEndian)
        {
            Array.Reverse(arr1);
        }

        int lenI = (int)BitConverter.ToUInt32(arr1, 0);
        byte[] InconNamebuf = new byte[lenI];
        InconNamebuf = data.Skip(4).Take(lenI).ToArray();
        string InconName = Encoding.UTF8.GetString(InconNamebuf);
        data = data.Skip(4 + lenI).ToArray();

        arr1[0] = data[0];
        arr1[1] = data[1];
        arr1[2] = data[2];
        arr1[3] = data[3];
        if (BitConverter.IsLittleEndian)
        {
            Array.Reverse(arr1);
        }

        int jinbi = (int)BitConverter.ToInt32(arr1, 0);

        UserDto dto = new UserDto(id, name, InconName, jinbi);
        return dto;

    }

    private static List<int> NewUpackListInt(byte[] data,int number)
    {
        List<int> readyUserIdList = new List<int>();

        for (int i=0;i<number;i++)
        {
            byte[] arr1 = new byte[4];

            arr1[0] = data[0];
            arr1[1] = data[1];
            arr1[2] = data[2];
            arr1[3] = data[3];
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(arr1);
            }

            int value = (int)BitConverter.ToUInt32(arr1, 0);
            data = data.Skip(4).ToArray();
            readyUserIdList.Add(value);

        }

        return readyUserIdList;
    }


    private static Dictionary<int, UserDto>  NewUpackuserIdUserDtoDic(byte[] data)
    {
        //用户map个数 + UserDto（value）
        Dictionary<int, UserDto> userIdUserDtoDic = new Dictionary<int, UserDto>();

        byte[] MapLen = new byte[4];

        MapLen[0] = data[0];
        MapLen[1] = data[1];
        MapLen[2] = data[2];
        MapLen[3] = data[3];
        if (BitConverter.IsLittleEndian)
        {
            Array.Reverse(MapLen);
        }

        int number = (int)BitConverter.ToUInt32(MapLen, 0);
        data = data.Skip(4).ToArray();

        for (int i=0;i<number;i++)
        {
            byte[] arr1 = new byte[4];

            arr1[0] = data[0];
            arr1[1] = data[1];
            arr1[2] = data[2];
            arr1[3] = data[3];
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(arr1);
            }

            int userdtolen = (int)BitConverter.ToUInt32(arr1, 0);
            data = data.Skip(4).Take(4 + userdtolen).ToArray();

            UserDto dto = NewUpackUserDto(data);

            userIdUserDtoDic.Add(dto.UserId, dto);
            

        }
        return userIdUserDtoDic;
    }

    /// <summary>
    /// 序列化
    /// @author wangfw
    /// @time 2019-4-2
    /// </summary>
    /// <param name="obj"></param>
    /// <returns></returns>
    public static byte[] NewEncodeObj(object obj)
    {


        using (MemoryStream ms = new MemoryStream())
        {
           
            NetMsg msg = obj as NetMsg;

            if (msg.value == null)   //只打包code
            {
                byte[] package = NewEncodeCode(msg.opCode, msg.subCode);


                return package;
            }
            Type t = msg.value.GetType();


            //打包消息体
            if (t == typeof(AccountDto))
            {
                //打包
                //List<byte> my = new List<byte>();

                AccountDto p = msg.value as AccountDto;

                //打包名字
                byte[] arr1 = Encoding.UTF8.GetBytes((string)p.userName);
                UInt16 len = (UInt16)arr1.Length;
                byte[] arr2 = BitConverter.GetBytes(len);
                if (BitConverter.IsLittleEndian)  //如果为xiaoduan
                {
                    Array.Reverse(arr2);
                }


                byte[] packageName = new byte[arr1.Length + arr2.Length];

                //先长度再内容
                arr2.CopyTo(packageName, 0);
                arr1.CopyTo(packageName, arr2.Length);

                //打包密码
                byte[] arr3 = Encoding.UTF8.GetBytes((string)p.password);
                UInt16 len2 = (UInt16)arr3.Length;
                byte[] arr4 = BitConverter.GetBytes(len2);
                if (BitConverter.IsLittleEndian)  //如果为xiaoduan
                {
                    Array.Reverse(arr4);
                }

                byte[] packagePas = new byte[arr3.Length + arr4.Length];

                //先长度再内容
                arr4.CopyTo(packagePas, 0);
                arr3.CopyTo(packagePas, arr4.Length);

                //打包消息体
                byte[] data = new byte[packageName.Length + packagePas.Length];
               
                packageName.CopyTo(data, 0);
                packagePas.CopyTo(data, packageName.Length);      

                //最终打包
                byte[] code = NewEncodeCode(msg.opCode, msg.subCode);
                byte[] package = new byte[code.Length + data.Length];
                code.CopyTo(package, 0);
                data.CopyTo(package, code.Length);


                return package;
            }
            else if(t == typeof(Int32))
            {
                byte[] buf = BitConverter.GetBytes((Int32)msg.value);
                if (BitConverter.IsLittleEndian)  //如果为xiaoduan
                {
                    Array.Reverse(buf);
                }

                byte[] codebuf = NewEncodeCode(msg.opCode, msg.subCode);
                byte[] package = new byte[buf.Length + codebuf.Length];
                codebuf.CopyTo(package, 0);
                buf.CopyTo(package, codebuf.Length);


                return package;
            }

            return null;
        }

    }






    //打包subcode和opcode
    /*王付伟 2019-4-3 测试 */
    private static byte[] NewEncodeCode(int Opcode,int Subcode)
    {
        //Int32 Opcode = msg.opCode;
       // Int32 Subcode = msg.subCode;
        byte[] arrOpcode = BitConverter.GetBytes((Int32)(Opcode));
   
        if (BitConverter.IsLittleEndian)  //如果为xiaoduan
        {
            Array.Reverse(arrOpcode);
        }

        byte[] arrSubcpde = BitConverter.GetBytes((Int32)(Subcode));

        if (BitConverter.IsLittleEndian)  //如果为xiaoduan
        {
            Array.Reverse(arrSubcpde);
        }

        byte[] packetCode = new byte[arrOpcode.Length + arrSubcpde.Length];
        arrOpcode.CopyTo(packetCode, 0);
        arrSubcpde.CopyTo(packetCode, arrOpcode.Length);

        return packetCode;

    }


}


