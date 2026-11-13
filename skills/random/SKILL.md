---
name: randomizer
description: Generates a random value for the requested type of object
---

# Functionality

Check the kind of object the user is asking about. If it is a number or a string, then use the `randomizer.go` program in the current directory to generate the random value. If the user is asking about a physical object, like for example a fruit, then generate a list of 10 different objects. Then generate a random number between 1 and 10 using the `randomizer.go` program, and return the object from the list at the index returned by the random number generator.

You can call the randomizer using this invocation:
`go run randomizer.go <lower> <upper>` where the <lower> and <upper> represent the lower and upper bounds of the range for the random number to be generated. For example, if you call randomizer using:
`go run randomizer 1 300` then you should get back a value between 1 and 300, for example 184.

## Examples
- Give me a random boy name
Assume that your request from the model for a list of 10 boy names gave this: Ram, Kishan, George, Mohammed, Shankar, Mahesh, Rahim, Soordas, Keshav, Madhav. Assume that your invocation of `randomizer.go` for a number between 1 and 10 gave the number 3. Then you would return the name "George" back to the user.

- Give me a random game
Assume that your request from the model for a list of 10 games gave this: cricket, baseball, football, hockey, chess, carrom, basketball, soccer, checkers and kabaddi. Assume that your invocation of `randomizer.go` for a number between 1 and 10 gave the number 7. Then you would return the game "basketball" back to the user.

## Guidelines
If you don't know the type of object that the user is asking about, it is perfectly fine to ask them for more details about their request. For example, a request for a random bffda9 should generate a response telling the user that you don't know anything that is called a "bffda9", and asking the user to clarify. Maybe the mistyped it, mis-spelt it, or are just testing you.

