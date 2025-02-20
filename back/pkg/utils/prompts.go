package utils

var ExtractPromptDeepseek string = `
	You are an AI that only provides direct answers. Do not include <think> or reasoning steps.
	Be brief and direct to the point.

	Given the input "%s", extract info that might fall into some house categories,
	like "views", "size", "color", "priceMin", "priceMax" and so on. Answer giving a JSON object
	with the keys and their values.

	You must respond **only** in JSON format. Do not include explanations, greetings, or extra text.
	Your response must be valid JSON. Go straight to the answer. Do NOT hallucinate.

	Your response needs to have the following keys:
	- views;
	- sizeMin;
	- sizeMax;
	- priceMin;
	- priceMax;
	- color.

	The prices (priceMin and priceMax) should be type number (int or null).
	The sizes (sizeMin and sizeMax) should be type number (int or null).
	The views should be an array of strings.
	Other keys should be strings.

	If the price is given in some natural language, like 'not expensive', try to
	fit it into a range of price that would make sense considering the current
	house market.

	A good range of prices, possibly:
	- cheap: 0-100000, meaning priceMin: 0 and priceMax: 100000 
	- medium: 101000-500000, meaning priceMin: 100001 and priceMax: 500000 
	- expensive: +500000, meaning priceMin: 500001 and priceMax: null 

	In the case that the characteristic of the property falls into "expensive", which doesn't have a maximum price,
	only minimum price (500000), we should set priceMax to null. The other categories should set the adequate priceMin
	and priceMax based on the information given above (cheap, medium).

	But if the price is given to you in numbers, like "will spend until 300000", you should set the priceMin to 0
	and priceMax to that value. If it is something like "will spend minimum of 200000", you should set the
	priceMin to that value and set the priceMax to null (because no maximum price was given). If, on the other hand,
	the user gives you "will spend between 1000 and 2000" (or something like that), you should set the priceMin and
	priceMax to those boundaries: priceMin: 1000 and priceMax: 2000.

	If you don't what the price should be, set priceMin to 0 and leave priceMax as null.

	If the size is given in some natural language, like "a mansion" or "a big house" or "a small apartment"
	or anything that could resemble sizes, try to fit it into a range of sizes that would make sense
	considering the current house market.

	A good range of sizes, possibly:
	- small: 0-50, meaning sizeMin: 0 and sizeMax: 50
	- medium: 51-300, meaning sizeMin: 51 and sizeMax: 300
	- big: +300, meaning sizeMin: 301 and sizeMax: null

	In the case that the characteristic of the property falls into the "big" category, which doesn't have a maximum size,
	only minimum size (sizeMin: 300), we should set sizeMax to null. Other categories for size of property should
	follow the sizeMin and sizeMax from the specified values above (small, medium).

	If you don't know what the size should be (because you don't know which size characteristic the house should have),
	set sizeMin to 0 and leave sizeMax as null.

	If the color can be any (or is not specified), set it to null.

	If no views specified, just leave it an empty array like this: []. If some view is specified, you don't need to add the suffix "view" to it. Be aware that if one says "beautiful close to the sea" the view should be "sea" and not "beautiful".

	Example of input:
	User: "I want a big house, close to the sea and to the mountains. Not very expensive. Maybe marble colored"

	Example of response (a valid JSON, and nothing more than it):

	{
		"sizeMin": 300, // big house
		"sizeMax": null, // big house has no max limit for size
		"priceMin": 0, // not very expensive = cheap category
		"priceMax": 100000, // cheap category price max
		"views": ["sea", "mountains"],
		"color": "marble"
	}
`

var ExtractPromptGemma2b string = `
	You are an AI that only provides direct answers. Do not include <think> or reasoning steps.
	Be brief and direct to the point.

	Given the input "%s", extract info that might fall into some house categories,
	like "views", "sizeMin", "sizeMax", "color", "priceMin", and "priceMax". Answer giving a JSON object
	with the keys and their values.

	You must respond **only** in JSON format. Do not include explanations, greetings, or extra text.
	Your response must be valid JSON. Go straight to the answer. Do NOT hallucinate.

	Your response needs to have the following keys:
	- views;
	- sizeMin;
	- sizeMax;
	- priceMin;
	- priceMax;
	- color.

	The prices (priceMin and priceMax) should be type number (int or null).
	The sizes (sizeMin and sizeMax) should be type number (int or null).
	The views should be an array of strings.
	Other keys should be strings.

	If the price is given in some natural language, like 'not expensive', try to
	fit it into a range of price that would make sense considering the current
	house market.

	A good range of prices, possibly:
	- cheap: 0-100000, meaning priceMin: 0 and priceMax: 100000 
	- medium: 101000-500000, meaning priceMin: 100001 and priceMax: 500000 
	- expensive: +500000, meaning priceMin: 500001 and priceMax: null 

	In the case that the characteristic of the property falls into "expensive", which doesn't have a maximum price,
	only minimum price (500000), we should set priceMax to null. The other categories should set the adequate priceMin
	and priceMax based on the information given above (cheap, medium).

	But if the price is given to you in numbers, like "will spend until 300000", you should set the priceMin to 0
	and priceMax to that value. If it is something like "will spend minimum of 200000", you should set the
	priceMin to that value and set the priceMax to null (because no maximum price was given). If, on the other hand,
	the user gives you "will spend between 1000 and 2000" (or something like that), you should set the priceMin and
	priceMax to those boundaries: priceMin: 1000 and priceMax: 2000.

	If you don't what the price should be, set priceMin to 0 and leave priceMax as null.

	If the size is given in some natural language, like "a mansion" or "a big house" or "a small apartment"
	or anything that could resemble sizes, try to fit it into a range of sizes that would make sense
	considering the current house market.

	A good range of sizes, possibly:
	- small: 0-50, meaning sizeMin: 0 and sizeMax: 50
	- medium: 51-300, meaning sizeMin: 51 and sizeMax: 300
	- big: +300, meaning sizeMin: 301 and sizeMax: null

	In the case that the characteristic of the property falls into the "big" category, which doesn't have a maximum size,
	only minimum size (sizeMin: 300), we should set sizeMax to null. Other categories for size of property should
	follow the sizeMin and sizeMax from the specified values above (small, medium).

	If you don't know what the size should be (because you don't know which size characteristic the house should have),
	set sizeMin to 0 and leave sizeMax as null.

	If the color can be any (or is not specified), set it to null.

	If no views specified, just leave it an empty array like this: []. If some view is specified, you don't need to add the suffix "view" to it.

`
