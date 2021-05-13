tokens:  36,000,000
blktime: 5
time:    1244160000(40*12*30*24*3600)
blks:    248832000(time/blktime)
AVGreward: 0.14467592592592593(tokens/blks)

[reward]
amplitude = random(0.001,AVGreward) // where amplitude is RANDOMLY generated every cycle changed.
cycle = random(100,3000) // where cycle amplitude is RANDOMLY generated every cycle finished.
radian = 2 * pi * step / cycle
reward = amplitude * sin(radian) + AVGreward