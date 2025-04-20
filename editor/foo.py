import sys
import struct

if __name__ == "__main__":
    datas = []

    while True:
        length_bytes = sys.stdin.buffer.read(4)
        if not length_bytes:
            break

        if len(length_bytes) < 4:
            print("incomplete length prefix", file=sys.stderr)
            break

        length = struct.unpack('>I', length_bytes)[0]

        data = sys.stdin.buffer.read(length)
        if len(data) < length:
            print("incomplete data", file=sys.stderr)
            break

        datas.append(data.decode('utf-8'))

    for item in reversed(datas):
        print(item)
