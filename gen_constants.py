import os

f = open("constants.go", "w")


def gen(s: str):
    f.write(s)


def genln(s: str):
    f.write(s + "\n")


def genint(n: int, prefix=True):
    n = "{0:#066x}".format(n)[2:]
    return (
        ("uint256.Int{" if prefix else "{")
        + "0x"
        + n[48:64]
        + ","
        + "0x"
        + n[32:48]
        + ","
        + "0x"
        + n[16:32]
        + ","
        + "0x"
        + n[0:16]
        + "}"
    )


genln("package main")

genln('import "github.com/holiman/uint256"')


LETTERS = ("A", "B", "C", "D", "E", "F", "G", "H", "I")

# generate squares
start = 51
genln("const (")
for i in range(10):
    gen(",".join(map(lambda x: x + str(i), LETTERS)))
    gen(" uint8 =")
    gen(",".join(map(lambda x: hex(x), range(start, start + 9))))
    gen("\n")
    start += 16
genln(")")

start = 51
genln("var Squares = [256]uint8{")
for i in range(10):
    gen(hex(start) + ":")
    genln(",".join(map(lambda x: hex(x), range(start, start + 9))) + ",")
    start += 16
genln("}")

# generate squares 180
start = 51
genln("var Squares180 = [256]uint8{")
for i in range(10):
    gen(hex(start) + ":")
    genln(",".join(map(lambda x: hex(x ^ 0xF0), range(start, start + 9))) + ",")
    start += 16
genln("}")

# generate bb squares
start = 51
genln("var (")
for i in range(10):
    j = start
    for a in LETTERS:
        gen("Bb" + a + str(i) + "=")
        genln(genint(1 << j).replace("0x0000000000000000", "0"))
        j += 1
    start += 16
genln(")")
# start = 51
# genln("var BbSquares = [256]uint256.Int{")
# for i in range(10):
#     gen(hex(start) + ":")
#     genln(",".join(map(lambda x: "Bb" + x+str(i), LETTERS)) + ",")
#     start += 16
# genln("}")

genln("var BbSquares = [256]uint256.Int{")
for i in range(128):
  gen(genint(1 << 2*i, False).replace("0x0000000000000000", "0")+",")
  genln(genint(1 << 2*i+1, False).replace("0x0000000000000000", "0")+",")
genln("}")

# generate files and ranks
genln("var (")
start = 3
for i in LETTERS:
    genln("File"+i+"="+genint(0x0001_0001_0001_0001_0001_0001_0001_0001_0001_0001_0001_0001_0001_0001_0001_0001<<start))
    start+=1
start = 3
for i in range(10):
    genln("Rank"+str(i)+"="+genint(0xffff << (16 * start)))
    start+=1
genln(")")

genln("var Files = [16]uint256.Int{")
genln("3:"+ ",".join(map(lambda x:"File"+x, LETTERS)) + ",")
genln("}")

genln("var Ranks = [16]uint256.Int{")
genln("3:"+ ",".join(map(lambda x:"Rank"+str(x), range(10))) + ",")
genln("}")


genln("var (")
genln("BbEmpty="+genint(0))
genln("BbAll="+genint(0xffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff))
genln("BbInBoard="+genint(0x0000_0000_0000_0ff8_0ff8_0ff8_0ff8_0ff8_0ff8_0ff8_0ff8_0ff8_0ff8_0000_0000_0000))
genln("BbRedSide="+genint(0x0000_0000_0000_0000_0000_0000_0000_0000_ffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff))
genln("BbBlackSide="+genint(0xffff_ffff_ffff_ffff_ffff_ffff_ffff_ffff_0000_0000_0000_0000_0000_0000_0000_0000))
genln(")")


f.close()
# format code
os.system("gofmt -w ./constants.go")
