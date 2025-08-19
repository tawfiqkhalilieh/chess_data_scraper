package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
    // This is the main function that starts the program, Uncomment the call for the function you want to run, if you want to run the entire program, uncomment all the calls. 

    // GamesCollector();

    FetchGameMovesFromJSON();
}

// Fetches 100k+ games of top players from chess.com and saves them to game_ids.json
func GamesCollector() {

    // 25 { "hikaru", "gmwso", "magnuscarlsen", "fabianocaruana", "lachesisq","lyonbeast", "hansontwitch", "vi_pranav", "VincentKeymer", "anishgiri", "gm_dmitrij", "firouzja2003", "denlaz", "0gzpanda", "javokhir_sindarov05", "rpragchess", "philippians46", "parhamov", "tradjabov", "raunaksadhwani2005", "spicycaterpillar", "chesswarrior7197", "mishanick", "andreikka", "chefshouse"} 
    // 8 {"grandelicious", "liemle", "wonderfultime", "liamputnam2008", "konavets", "shield12", "tptagain", "xiaotong2008" }
    // 8 "arseniy_nesterov",  "danielnaroditsky", "macho_2006", "robert_chessmood", "hovik_hayrapetyan", "sanan_sjugirov", "danieldardha2005", "viditchess"

    var players = [3] string {
        "santoshgupta", "mariogiri", "gmbenjaminbok"}

    for _, name := range players {
        var i int = 0;
        var current_pages int = 0;
        var game_ids[] string

        for (current_pages == 0 || i != current_pages) {
            var url string = "https://www.chess.com/callback/games/extended-archive?locale=en&username=" + name + "&page=" + strconv.Itoa(current_pages) + "&gameResult=won&gameTypeslive%5B%5D=blitz&gameTypeslive%5B%5D=chess&rated=rated&timeSort=desc&location=live&opponentTitle=GM&result=won"

            resp, err := http.Get(url)
            if err != nil {
                log.Fatalln(err)
            }

            if (resp.StatusCode == 200) {
                defer resp.Body.Close() // Ensure the response body is closed

                // Read the response body
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                    fmt.Printf("Error reading response body: %v\n", err)
                    return
                }


                var result map[string] interface {}
                json.Unmarshal(body, & result)

                current_pages += 1;

                if (i == 0) {
                    // === totalPages: result["meta"]["totalPages"] ===
                    totalPages := -1
                    if metaRaw,
                    ok := result["meta"];ok && metaRaw != nil {
                        if meta, ok := metaRaw.(map[string] interface {});
                        ok {
                            if tpRaw, exists := meta["totalPages"];
                            exists && tpRaw != nil {
                                switch v := tpRaw.(type) {
                                    case float64:
                                        totalPages = int(v)
                                    case string:
                                        if n, err := strconv.Atoi(v);
                                        err == nil {
                                            totalPages = n
                                        }
                                    case json.Number: // only if you used Decoder.UseNumber()
                                        if n, err := v.Int64();
                                        err == nil {
                                            totalPages = int(n)
                                        }
                                    default:
                                        // unexpected type
                                }
                            }
                        }
                    }
                    i = totalPages
                }

                // === IDs: result["data"][i]["id"] ===
                if dataRaw, ok := result["data"].([] interface {});
                ok {
                    for i, item := range dataRaw {
                        if obj, ok := item.(map[string] interface {});
                        ok {

                            // id might be string or number depending on API
                            if idRaw, ok := obj["id"];
                            ok && idRaw != nil {
                                switch id := idRaw.(type) {
                                    case string:
                                        // fmt.Println(id + "\n")
                                        game_ids = append(game_ids, id)
                                    case float64:
                                        // if IDs are numeric
                                        game_ids = append(game_ids, strconv.FormatFloat(id, 'f', -1, 64))
                                            // fmt.Printf(strconv.FormatFloat(id, 'f', -1, 64) + "\n")
                                    default:
                                        fmt.Printf("%d \n", i)
                                }
                            } else {
                                // fmt.Printf("%d \n", i)
                            }
                        }
                    }
                } else {
                    // fmt.Println("data not found or not an array")
                }
            } else {

                if (resp.StatusCode == 429) {
                    time.Sleep(30 * time.Second) // Pauses execution for 30 seconds 
                } else {
                    addToJSON(name, game_ids)
                    fmt.Println("Error fetching data:", resp.StatusCode)
                    game_ids = [] string {}
                    i = 0
                    current_pages = 0

                }

            }
        }

        if (len(game_ids) > 0) {
            addToJSON(name, game_ids)
        }
    }
}

func printSlice(s[] string) {
    fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

// addToJSON adds a key with a string array to data.json
func addToJSON(key string, values[] string) error {
    filename := "game_ids.json"
    data := make(map[string] interface {})

    // If file exists, load it
    if _,
    err := os.Stat(filename);err == nil {
        fileBytes, err := ioutil.ReadFile(filename)
        if err != nil {
            return fmt.Errorf("failed to read file: %v", err)
        }
        if len(fileBytes) > 0 {
            if err := json.Unmarshal(fileBytes, & data);
            err != nil {
                return fmt.Errorf("failed to parse JSON: %v", err)
            }
        }
    }

    // Add or update key
    data[key] = values

    // Marshal with indentation for readability
    output,
    err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to encode JSON: %v", err)
    }

    // Save back to file
    if err := ioutil.WriteFile(filename, output, 0644);err != nil {
        return fmt.Errorf("failed to write file: %v", err)
    }

    return nil
}

type Move struct {
    From string
    To string
    Drop * string
    Promotion * string
}

func indexOf(s string, c rune) int {
    for i, ch := range s {
        if ch == c {
            return i
        }
    }
    return -1
}

func decodeTCN(n string)[] Move {
    tcnChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!?{~}(^)[_]@#$,./&-*++="
    pieceChars := "qnrbkp"

        c := [] Move {}
    w := len(n)

    for i := 0;i < w;i += 2 {
        var u Move
        o := indexOf(tcnChars, rune(n[i]))
        s := indexOf(tcnChars, rune(n[i + 1]))

        if s > 63 {
            promIndex := int(math.Floor(float64(s - 64) / 3))
            promotionChar := string(pieceChars[promIndex])
            u.Promotion = & promotionChar
            // Calculate s for the "to" square
            if o < 16 {
                s = o + (-8) + ((s - 1) % 3) - 1
            } else {
                s = o + 8 + ((s - 1) % 3) - 1
            }
        }

        if o > 75 {
            dropChar := string(pieceChars[o - 79])
            u.Drop = & dropChar
        } else {
            fromSquare := string(tcnChars[o % 8]) + fmt.Sprintf("%d", int(math.Floor(float64(o) / 8)) + 1)
            u.From = fromSquare
        }

        toSquare := string(tcnChars[s % 8]) + fmt.Sprintf("%d", int(math.Floor(float64(s) / 8)) + 1)
        u.To = toSquare

        c = append(c, u)
    }

    return c
}

// Struct to match the JSON structure
type GameResponse struct {
    Game struct {
        MoveList string `json:"moveList"`
        MoveTimestamps string `json:"moveTimestamps"`
    } `json:"game"`
}

func GetChessGameData(id string)(string, string, error) {
    url := fmt.Sprintf("https://www.chess.com/callback/live/game/%s", id)

        resp,
    err := http.Get(url)
    if err != nil {
        logError(id) // Log the error to a file
        time.Sleep(120 * time.Second) // Pauses execution for 120 seconds 
        return "", "", nil
        // return "", "", fmt.Errorf("failed to send GET request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        logError(id) // Log the error to a file
        time.Sleep(120 * time.Second) // Pauses execution for 120 seconds 
        return "", "", nil
        // return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    body,
    err := io.ReadAll(resp.Body)
    if err != nil {
        logError(id) // Log the error to a file
        time.Sleep(120 * time.Second) // Pauses execution for 120 seconds 
        return "", "", nil
        // return "", "", fmt.Errorf("failed to read response body: %w", err)
    }

    var gameData GameResponse
    if err := json.Unmarshal(body, & gameData);err != nil {
        logError(id) // Log the error to a file
        time.Sleep(120 * time.Second) // Pauses execution for 120 seconds 
        return "", "", nil
        // return "", "", fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    return gameData.Game.MoveList,
    gameData.Game.MoveTimestamps,
    err
}

func FetchGameMovesFromJSON() {
    // Open the JSON file
    file, err := os.Open("game_ids.json")
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer file.Close()

    // Read the file
    data, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatalf("Error reading file: %v", err)
    }

    // Parse JSON into a generic map
    var gameData map[string][] string
    err = json.Unmarshal(data, & gameData)
    if err != nil {
        log.Fatalf("Error parsing JSON: %v", err)
    }

    var game_count int = 0;
    var multi_file_count int = 0;

    // File path
    filePath := "output_data/game_information" + strconv.Itoa(multi_file_count) + ".json"
    // Step 1: Load existing JSON or create new map
    gameInfo := make(map[string] map[string] interface {})

    if _, err := os.Stat(filePath);
    err == nil {
        // File exists, read and unmarshal
        data, err := ioutil.ReadFile(filePath)
        if err != nil {
            log.Fatalf("Error reading file: %v", err)
        }
        if len(data) > 0 {
            err = json.Unmarshal(data, & gameInfo)
            if err != nil {
                log.Fatalf("Error parsing JSON: %v", err)
            }
        }
    }

    // Loop over the keys and arrays
    for key, values := range gameData {
        fmt.Printf("%s\n", key)

        /*
        var emeregency string = " gm_dmitrijliamputnam2008magnuscarlsenmishanickrpragchessshield12xiaotong2008denlazfabianocaruanahovik_hayrapetyanjavokhir_sindarov05spicycaterpillarwonderfultimegrandelicioushikarukonavets"
        if (strings.Contains(emeregency, key)) {
            fmt.Printf("key: %s\n", key)
            continue
        } else {
            fmt.Printf("Key: %s\n", key)
        }
            */
        
        for _, value := range values {
            var resultMovesArray = [] string {}
            var MoveTimestampsArray[] string

            var whiteMoveTimestampsArray[] string
            var blackMoveTimestampsArray[] string


            // fmt.Printf("  Value: %s\n", value)

            moveList, moveTimestamps, err := GetChessGameData(value) // replace with a real game ID
            if err != nil {
                fmt.Println("Error:", err)
                return
            }

            moveListArray := decodeTCN(moveList)

            for _, move := range moveListArray {
                resultMovesArray = append(resultMovesArray, move.From + move.To)
            }

            MoveTimestampsArray = strings.Split(moveTimestamps, ",")

            for i, stamp := range MoveTimestampsArray {
                if i % 2 == 1 {
                    whiteMoveTimestampsArray = append(whiteMoveTimestampsArray, stamp)
                } else {
                    blackMoveTimestampsArray = append(blackMoveTimestampsArray, stamp)
                }
            }

            

            // Step 2: Add/overwrite this game's data
            gameInfo[value] = map[string] interface {} {
                "moveListArray": resultMovesArray,
                "whiteMoveTimestampsArray": whiteMoveTimestampsArray,
                "blackMoveTimestampsArray": blackMoveTimestampsArray,
            }

            game_count += 1;
            
            // Step 3: Write back to file
            output, err := json.MarshalIndent(gameInfo, "", "  ")

            if err != nil { 
                log.Fatalf("Error encoding JSON: %v", err)
            }

            err = ioutil.WriteFile(filePath, output, 0644)
            if err != nil {
                log.Fatalf("Error writing file: %v", err)
            }

            if ( game_count > 5000) {
                multi_file_count += 1;
                game_count = 0;
                filePath = "output_data/game_information" + strconv.Itoa(multi_file_count) + ".json"
            }

            time.Sleep(250 * time.Millisecond) // Pauses execution for 30 seconds 
                // fmt.Println("Game information saved successfully.")
        }

        
    }
}


func logError(game_id string) {
    filename := "errors.txt"
    dataToAppend := game_id + "\n"

    // Open the file in append mode, create if it doesn't exist, and set write-only permissions
        f,
    err := os.OpenFile(filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer f.Close() // Ensure the file is closed

    // Write the data to the file
    if _,
    err := f.WriteString(dataToAppend);err != nil {
        log.Fatalf("failed writing to file: %s", err)
    }

}