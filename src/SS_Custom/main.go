package SS_Main

/*
        +--------------------- IMPORTING ---------------------+
	|          To import other things in the same         |
	|            "/src/" folder simply include            |
	+-----------------------------------------------------+
	|                "main/src/SS_Custom/"                |
	+-----------------------------------------------------+

	+------------- CUSTOM AUDIO IN CALLBACK --------------+
        |   If you want a certain audio to be played in a     |
	|        a custom callback you can import             |
	|     "main/src/AudioPlayer"' and call the New()      |
	|   Function with the path as a parameter. It will    |
	|   return *AudioPlayer and/or an error. Check for    |
	|  errors and you can call the *AudioPlayer.Play()    |
	|             function to play the sound              |
	+-----------------------------------------------------+
*/

func Custom_Main(EventTitle string, EventTimestamp int64) {
	// ...
}
