﻿using System;
using System.Collections.Generic;
using OpenDiablo2.Common.Enums;
using OpenDiablo2.Common.Models.Mobs;

namespace OpenDiablo2.Common.Models
{
    public sealed class PlayerLocationDetails
    {
        public int PlayerId { get; set; }
        public float PlayerX { get; set; }
        public float PlayerY { get; set; }
        public eMovementType MovementType { get; set; }
        public int MovementDirection { get; set; }
        public float MovementSpeed { get; set; }
        // TODO: They may not be on the same 'anchor map'...

        public byte[] GetBytes()
        {
            var result = new List<byte>();
            result.AddRange(BitConverter.GetBytes((Int32)PlayerId));
            result.AddRange(BitConverter.GetBytes((float)PlayerX));
            result.AddRange(BitConverter.GetBytes((float)PlayerY));
            result.AddRange(BitConverter.GetBytes((Int32)MovementDirection));
            result.AddRange(BitConverter.GetBytes((byte)MovementType));
            result.AddRange(BitConverter.GetBytes((float)MovementSpeed));
            return result.ToArray();
        }

        public static PlayerLocationDetails FromBytes(byte[] data, int offset = 0)
        {
            var result = new PlayerLocationDetails
            {
                PlayerId = BitConverter.ToInt32(data, offset + 0),
                PlayerX = BitConverter.ToSingle(data, offset + 4),
                PlayerY = BitConverter.ToSingle(data, offset + 8),
                MovementDirection = BitConverter.ToInt32(data, offset + 12),
                MovementType = (eMovementType)data[offset + 16],
                MovementSpeed = BitConverter.ToSingle(data, offset + 18)
            };
            return result;
        }
        public static int SizeInBytes => 22;
    }

    public static class PlayerLocationDetailsExtensions
    {
        public static PlayerLocationDetails ToPlayerLocationDetails(this PlayerState source)
        {
            var result = new PlayerLocationDetails
            {
                PlayerId = source.Id,
                PlayerX = source.GetPosition().X,
                PlayerY = source.GetPosition().Y,
                MovementType = source.MovementType,
                MovementDirection = source.MovementDirection,
                MovementSpeed = (float)(source.MovementType == eMovementType.Running ? source.GetRunVelocity() : source.GetWalkVeloicty()) / 4f
            };
            return result;
        }
    }
}