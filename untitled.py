

funcs = [lambda x: x+i for i in range(10)]

funcs = []
for i in range(10):
	funcs.append(lambda x: x+i)

print sum(f(1) for f in funcs)


def test():
	funcs = []
	for i in range(10):
		def closure(i):
			def inside(x):
				return x+i
			return inside
		funcs.append(closure)


if __name__ == '__main__':
	test()