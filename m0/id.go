package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	uuid "github.com/satori/go.uuid"
)

var id int64

func init() {
	rand.Seed(time.Now().UnixNano())

	id = rand.Int63n(100000000) + 100000000
}

func GetID(name ...string) string {
	if len(name) == 0 {
		return GetRandomName() // 默认为 name
	}

	switch name[0] {
	case "increase":
		return GetIncreaseID()
	case "uuid":
		return GetUUID()
	case "random":
		return GetRandomStr(8)
	case "name":
		return GetRandomName()
	default:
		return GetRandomStr(8)
	}
}

func GetIncreaseID() string {
	return strconv.Itoa(int(atomic.AddInt64(&id, 1)))
}

func GetUUID() string {
	return uuid.NewV4().String()
}

var ss = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

func GetRandomStr(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = ss[rand.Intn(len(ss))]
	}
	return string(b)
}

var randNameSeed = []string{"Thunder", "Martha", "Edric", "Derwin", "Shamus", "Joshua", "Church", "Edgar", "Renata", "Hannah", "Dillon", "Nerita", "Foster", "Rodney", "Travis", "Ryan", "Ardent", "Theodora", "Sean", "Bobbie", "Irving", "Joanna", "Helena", "Wilbur", "Louisa", "Emmanuel", "Karena", "Lois", "Oprah", "Marlon", "Healthy", "Fiona", "Herbert", "Lucinda", "Dominica", "Trevor", "Nessa", "Guinevere", "Heathcliff", "Tina", "Frances", "Adelaide", "Awe-Inspiring", "Victor", "Tuesday", "Halsey", "Eugene", "Edwin", "Elaine", "Jed", "Victorious", "Nola", "Prunella", "Valley", "Julie", "Earthy", "Maiden", "Alice", "Estra", "Harriet", "Elfin", "Lillian", "Magdalene", "Ezra", "Lombard", "Matthew", "Gladys", "Luminous", "Egbert", "Gideon", "Igor", "Frederick", "Theobold", "Bernice", "Landon", "Maude", "Gilroy", "Eldwin", "Ebenezer", "Forest", "Handsome", "Primavera", "Belinda", "Teresa", "Guide", "Sheila", "Jasmine", "Briana", "Attendant", "Lame", "Emrick", "Shana", "Robert", "Lea", "Crown", "Endurance", "Trustworthy", "Elena", "Sharon", "Kelsey", "Justine", "Imogen", "Roxanne", "Torrent", "Ralph", "Louise", "Percival", "Tabitha", "Pierce", "Elmer", "Joey", "Hilda", "Kara", "Compassionate", "Dexterous", "Gale", "Kirk", "Owner", "Lee", "Serpent", "Quincy", "Jacob", "Simon", "Dirk", "Melville", "Beata", "Strawberry", "Hugo", "Griswald", "Edmund", "Ariana", "Tamara", "Stacy", "Simone", "Tony", "Norine", "Julia", "Ferguson", "Harlan", "Willa", "Weaver", "Timothea", "Merle", "Ross", "Octavia", "Sybil", "Admirable", "Eudora", "Patty", "Rolf", "Barbara", "Leo", "Tara", "Douglas", "Ruth", "Livia", "Listener", "Haley", "Great", "Ivan", "Davin", "Lorena", "Lolita", "Winston", "Sharp", "Nerissa", "Lars", "Abigail", "Kirstyn", "Britney", "Kirby", "John", "Tasha", "Noble", "Loveable", "Quinella", "Garrick", "Leroy", "Perry", "Delight", "Keegan", "Marie", "Sidney", "Linette", "Beatrice", "Melody", "Sherwin", "Morris", "Strong", "Dutiful", "Idelle", "Jessica", "Jeanne", "Jemima", "Cunning", "Martina", "Kingsley", "Majestic", "Flower-Like", "Ulrica", "Questa", "David", "Hunter", "Bridget", "Dean", "Nobleman", "Beverly", "Powerful", "Penelope", "Declan", "Daley", "Anastasia", "Alma", "Famous", "Yvonne", "Juliet", "Bright", "Elvira", "Beloved", "Peaceful", "Vandal", "Wise", "Desmond", "Medwin", "Trix", "Edna", "Peter", "Blueberry", "Rosalind", "Kelsey", "Marian", "Zera", "Flourishing", "Holly", "Long-Beard", "Maggie", "Black", "Ethan", "Lola", "Sterling", "Belle", "Peyton", "Hortense", "Hector", "Dominique", "Rupert", "Morgan", "Anthea", "Maurice", "Maia", "Thea", "Sarah", "Hamlin", "Rory", "Lulu", "Dale", "Gardner", "Eleanor", "Howard", "Industrious", "Kit", "Landry", "Star", "Gilda", "Lambert", "Polly", "Peaceful", "Mildred", "Dorian", "Derek", "Sigmund", "Efrain", "Nightingale", "April", "Peacemaker", "Earth", "Ann", "Naomi", "Fenton", "Kelvin", "Astrid", "Quinn", "Glenn", "Hetty", "Island", "Erik", "Mavis", "Mercy", "Graham", "Imagine", "Effie", "Joan", "Happy", "Becky", "Wynne", "Faith", "Wyatt", "Sebastian", "Lawyer", "Oliver", "United", "Natalie", "Diane", "Leslie", "Gardener", "Guy", "Royal", "Uriah", "Kacey", "Gerard", "Paul", "Earl", "Owen", "Stream", "Caretaker", "King", "Godfrey", "Nicolette", "Leith", "Wolf", "Gertrude", "Daphne", "Orlena", "Gardener", "Luke", "Miriam", "Warlike", "Tammy", "Shelley", "Magnus", "Esmond", "Willette", "Salt", "Mary", "Ronald", "Laurence", "Sparrow", "Faye", "Harley", "Ursa", "Wendy", "Vivian", "James", "Travers", "Hubert", "Roswell", "Zoe", "Female", "Lion-like", "Ambitious", "Brooke", "Washington", "Seeds", "Stephen", "Keith", "Pure", "Prudent", "Soldier", "Victoria", "Emmett", "Alda", "Odette", "Kirsten", "Des", "Blessed", "Erin", "Leonard", "Loralie", "Gabrielle", "Ellen", "Geraldine", "Danielle", "Honor", "Lester", "Larissa", "Veleda", "Farmer", "Blooming", "Phoebe", "Francis", "Mariner", "Yvette", "Warrior", "Jarvis", "Kenyon", "Gillian", "Cheerful", "Power", "Falkner", "Ivory", "Fleming", "Salena", "Matilda", "Humphrey", "Jason", "Theodore", "Maxine", "Nadine", "Logan", "Ernestine", "Hope", "Liza", "Zelene", "Wilona", "Frank", "Thelma", "Genevieve", "Freeman", "Yolanda", "Fair-Haired", "Randolph", "Edan", "Jessie", "Mountain", "Precious", "Ely", "Ives", "Delmar", "Valerie", "Flora", "Aimee", "Renfred", "Herdsman", "Kayla", "Kayleigh", "Hugh", "Silver", "Wallace", "Driscoll", "Ancestress", "Russell", "Evelyn", "Donna", "Lorraine", "Kilian", "Jonathan", "Lynn", "Zachary", "Valiant", "Mark", "Quade", "Ridley", "Oriel", "Molly", "Estelle", "Vaughan", "Harley", "Meadow", "Garth", "Norma", "Konrad", "Halbert", "Tatum", "Melvina", "Ingrid", "Merlin", "Steadfast", "Megan", "Ophelia", "Robust", "Rhoda", "Roger", "Lucille", "Farrell", "Fiery", "Timekeeper", "Damon", "Elizabeth", "Samson", "Darell", "Flame", "Melinda", "Elsie", "William", "Trista", "Ellery", "Jasper", "Graceful", "Kate", "Beautiful", "Dragon", "Aileen", "Francesca", "Kingly", "Eldon", "Garret", "Stefan", "Madge", "Solomon", "Kurt", "Rita", "Industrious", "Zebediah", "Lucy", "Talia", "Darlene", "Giles", "Evan", "Rosa", "Philbert", "Beth", "Lucas", "Vigour", "Tristan", "Maura", "Leah", "Melissa", "Peace", "Iris", "Eloise", "Salome", "Peg", "Humble", "Gaye", "Dark-Haired", "Joe", "Sibyl", "Ward", "Flower", "Erika", "Luciana", "Norris", "Willow", "Fresh", "Kendrick", "Fabian", "Gwen", "Blind", "Brenda", "Willard", "Renee", "Marissa", "Dexter", "Herman", "Marta", "Kenneth", "Agatha", "Youthful", "Katherine", "Tiffany", "Geneva", "Rex", "Bettina", "Meris", "Nigel", "Ivy", "Quenby", "Lighthearted", "Dalton", "Woodsman", "Jennifer", "Egan", "Tanya", "Gavin", "Exalted", "Free", "Jack", "Beatrix", "Emery", "Joyce", "Interpreter", "Elise", "Elga", "Blythe", "Gazelle", "Nathaniel", "Serene", "Eagle-Eyed", "Gwendolyn", "Riley", "Tower", "Virtuous", "Bernadette", "Maxwell", "Gifted", "Half-Dane", "Marnia", "Norman", "Larina", "Lilly", "Magda", "Darcy", "Edward", "Mabel", "Zane", "Edmond", "Moorish", "Lively", "Falcon", "Rebecca", "Strange", "Samuel", "Rejoicing", "Grateful", "Amiable", "Jewel", "Fergal", "Angelica", "Blackbird", "Linda", "Page", "Hadley", "Dennis", "Noelle", "Grey", "Forbes", "Kendra", "Otis", "Tilda", "Archer", "Counsellor", "Kay", "Audrey", "Gilbert", "Priscilla", "Annabelle", "Lloyd", "Kane", "Leticia", "Harrison", "Olivia", "Red", "Anita", "Keene", "Hal", "Hall", "Hardy", "Griswold", "Dark", "Valda", "Kerry", "Kendall", "Vance", "Mark", "Moira", "Tyler", "Queenie", "Blanche", "Rhea", "Elvis", "Monica", "Lamont", "Free", "Eli", "Quentin", "Honour", "Whitney", "Holy", "Mandy", "Stewart", "Spirited", "Violet", "Georgette", "Fedora", "Brina", "Emily", "Sherard", "Walter", "Song-Thrush", "Lee", "Rebellious", "Godwin", "Milburn", "Drucilla", "Winifred", "Denise", "Lara", "Jocelyn", "Searcher", "Dylan", "Rhett", "Eugenia", "Unity", "Winona", "Cheerful", "Gerald", "Eliza", "Dependable", "Wealthy", "Fitzgerald", "Laverna", "Opal", "Sherlock", "Swift", "Zachariah", "Bess", "Jeffrey", "Gerret", "Quimby", "Oscar", "Vita", "Praised", "Trent", "Laughter", "Lowell", "Victorious", "Warren", "Keely", "Marcia", "Eva", "Wolf", "Will", "Grace", "Fairy", "Drake", "Mighty", "Rachel", "Henrietta", "Alina", "Beneficient", "Olin", "Olga", "Dwayne", "Strength", "Vivianne", "Arleen", "Everett", "Helen", "Paula", "Ken", "Thora", "Egil", "Lacey", "Pretty", "Elroy", "Gentle", "Plains", "Verda", "Olive", "Queenly", "Reginald", "Sparkling", "Conqueror", "Sandra", "God-like", "Myrtle", "Kevin", "Berta", "Egerton", "Dark", "Troy", "Sibley", "Lyndon", "Duncan", "Bertina", "Luna", "Harry", "Fiery", "Seth", "Georgia", "Bertha", "Primrose", "Wyman", "Articulate", "Lisa", "Kent", "Agnes", "Alarice", "Nydia", "Brittany", "Simona", "Meadow", "Edana", "Horace", "Kathleen", "Life", "Light", "Ian", "Sadie", "Maureen", "Philomena", "Vincent", "Patriotic", "Echo", "Marilyn", "Nathan", "Wesley", "Bonnie", "Elliott", "Misty", "Tobias", "Kathy", "Thresher", "Nessia", "Servant", "Famous", "Laurel", "Kimball", "Hanley", "Raymond", "Laura", "Willis", "Immortal", "Quenna", "Tracy", "Una", "Eileen", "Teri", "Robin", "Lorelei", "Forrest", "Maria", "Nell", "Lizzie", "Neal", "Walton", "Jimmy", "Farrah", "Gale", "Kerwin", "Fairfax", "Jesse", "Cub", "Ula", "Queen", "Small", "Sabrina", "Armed", "Veronica", "Shining", "Janice", "Roy", "Beauty", "Eric", "Gregory", "Harland", "Harold", "Eve", "Nursing", "Pledge", "Wanderer", "Durwin", "Martin", "Harmony", "Stranger", "Heroine", "Lilah", "Firm", "Glorious", "Elbert", "Zea", "Supplanter", "White", "Fawn", "Silas", "Nourishing", "Lewis", "Titus", "Malcolm", "Murray", "Henry", "Red", "Phyllis", "Katrina", "Steward", "Patricia", "Joy", "Frederica", "Felicia", "Loyal", "Prudence", "Irvin", "Quinby", "Bound", "Dixon", "Nicole", "Wonderful", "Patience", "Udele", "Edwina", "Heather", "Juliana", "Paxton", "Montgomery", "Fox", "Amber", "Orva", "Sea", "Shawn", "Deirdre", "Ada", "Hale", "Orlantha", "Fara", "George", "Jacqueline", "Ethel", "Nimble", "Lane", "Marc", "Grain", "Kilby", "Vanessa", "Floyd", "Grant", "Great", "Jill", "Luther", "Tess", "Max", "Amanda", "Butterfly", "Robin", "Millicent", "Imogene", "Unwin", "Eunice", "Rowena", "Darian", "Noel", "Miranda", "Scarlett", "Red-Haired", "Myrrh", "Lucia", "Pleasure", "Esmeralda", "Eaton", "Beryl", "Wayne", "Phineas", "Glynnis", "Edith", "Peggy", "Just", "Andrea", "Commander", "Holly", "Anne", "Ramsey", "Myra", "Iver", "Philippa", "Jillian", "Heath", "Harvester", "Dwight", "Respected", "Felix", "Samantha", "Relic", "Constant", "Gifford", "Margot", "Quintana", "Solitary", "Roberta", "Emma", "Adrienne", "Guardian", "Alanna", "Grover", "Wide", "Wylie", "Sapphire", "Fannie", "Mora", "Eddie", "Zelda", "Leanne", "Melvin", "Hadden", "Bella", "Annette", "Vernon", "Wealthy", "Michael", "Pearl", "Angelic", "Judith", "Helpful", "Ruby", "Dora", "Errol", "Seaman", "Rose", "Udolf", "Strong", "Toby", "Elton", "Rosemary", "Hayley", "Ramona", "Marcus", "Quintessa", "Zebadiah", "Morgan", "Nora", "Dudley", "Shannon", "Fergus", "Philip", "Lovely", "Grayson", "Wilda", "Warrior", "Georgiana", "Jade", "Elijah", "Weary"}

func GetRandomName() string {
	n := rand.Intn(10000)
	name := randNameSeed[rand.Intn(len(randNameSeed))]

	return fmt.Sprintf("%s_%d", name, n)
}
