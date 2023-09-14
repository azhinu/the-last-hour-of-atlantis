package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/pterm/pterm"
)

// Clear screen sequence
// os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
type Status struct {
	time int
	intro bool
	familySaved bool
	sailorFamilySaved bool
	artistsSaved bool
	scientistsSaved bool
	portVisited bool
	homeVisited bool
	castleVisited bool
}
type Location struct {
	port bool
	castle bool
	home bool
	artists bool
	scientists bool
}

// Start text ares
var area, _ = pterm.DefaultArea.WithFullscreen().WithCenter().Start()

// Game reset status
var gameEnd = false
// Levels status
var status = Status {
	time: 0,
	intro: false,
	familySaved: false,
	sailorFamilySaved: false,
	artistsSaved: false,
	scientistsSaved: false,
	portVisited: false,
	homeVisited: false,
	castleVisited: false,
}
var in = Location {
	port: false,
	castle: false,
	home: false,
	artists: false,
	scientists: false,
}

func main (){
	SetupCloseHandler()
	banner()
	for !gameEnd {
		if !status.intro {
			intro()
		}
		moveOptions := []string{"Отправиться в порт.", "Отправиться к городской площади.", "Отправиться в гильдию художников.", "Отправиться в резиденцию научного общества.", "Отправиться во дворец.", "Отправиться к семье."}
		result, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Что будете делать?").WithMaxHeight(7).WithOptions(moveOptions).Show()
		updateMessage("Вы отправились " + result + ".")
		switch result {
			case "Отправиться в порт.":
				portLoc()
			case "Отправиться на городскую площадь.":
				squareLoc()
			case "Отправиться в гильдию художников.":
				artistsLoc()
			case "Отправиться в резиденцию научного общества.":
				scientistsLoc()
			case "Отправиться во дворец.":
				castleLoc()
			case "Отправиться к семье.":
				homeLoc()
		}
		endOfLife()
	}
	area.Stop()
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pterm.ThemeDefault.InfoMessageStyle.Println("Exiting...")
		os.Exit(0)
	}()
}

func keyPressListener(startGame *bool) {
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	_, _, err = keyboard.GetSingleKey()
	if err != nil {
		log.Fatal(err)
	}
	*startGame = true
}

func updateMessage(message string) {
	area.Update(pterm.DefaultBox.Sprint(message))
}

func banner() {
	startGame := false
	go keyPressListener(&startGame) // Start listening any key press
	
	label := pterm.DefaultBasicText.Sprint("     ,--,--'.        .          .   .\n     `- |   |-. ,-.  |  ,-. ,-. |-  |-. ,-. . . ,-.  ,-. ,\"\n      , |   | | |-'  |  ,-| `-. |   | | | | | | |    | | |-\n      `-'   ' ' `-'  `' `-^ `-' `'  ' ' `-' `-^ '    `-' |\n                ,.   .  .          .      .              '\n               / |   |- |  ,-. ,-. |- . ,-| ,-.\n              /~~|-. |  |  ,-| | | |  | | | ,-|\n            ,'   `-' `' `' `-^ ' ' `' ' `-^ `-^\n")
	
	text := pterm.DefaultBasicText.Sprint("\n                    Press any key to continue")

	// Dot animation
	dots := ""
	for !startGame {
		switch dots{
			case ".":
				dots = ".."
			case "..":
				dots = "..."
			default:
				dots = "."
		}
		// Update area text
		area.Update(pterm.DefaultBox.WithRightPadding(7).WithBottomPadding(1).Sprint(label, text, dots))
		time.Sleep(time.Second)
	}
}

// Game over
func reset() {
	options := []string{"Начать заново.", "Выйти из игры."}
	result, _ := pterm.DefaultInteractiveSelect.
								WithDefaultText("Конец. Попробовать ещё раз?").
								WithOptions(options).Show()
	switch result {
		case "Начать заново.":
			status = Status {
				time: 0,
				intro: false,
				familySaved: false,
				sailorFamilySaved: false,
				artistsSaved: false,
				scientistsSaved: false,
				portVisited: false,
				homeVisited: false,
				castleVisited: false,
			}
			in = Location {
				port: false,
				castle: false,
				home: false,
				artists: false,
				scientists: false,
			}
		case "Выйти из игры.":
			status.time = 0
			gameEnd = true
		}
}

func resetLoc() {
	in = Location {
		port: false,
		castle: false,
		home: false,
		artists: false,
		scientists: false,
	}
}

func invalidLoc() {
	updateMessage("Вы пытались прийти в место, где вы уже есть. Вы заблудились в лесу сомнений.\nПытаясь выбраться из него, вы попали под текстуры, которых здесь нет.")
	reset()
}
// Ask to comfirm message
func nextMessage() {
	pterm.DefaultInteractiveSelect.WithDefaultText("").WithOptions([]string{"Ок"}).Show()
}

func intro() {
	area.Update(pterm.DefaultBox.Sprint("Вы придворный советник. К вам приходит пророк. Он показывает вам\nвидение. Оно подтверждает самые страшные опасения: Атлантида падёт.\nОстался всего час до гибели целого государства.\nВы отправляете к правителю. Он поручает вам миссию по спасению народа.\nТеперь вам решать, кто спасётся."))
	status.intro = true
}

func portLoc() {
	if in.port {
		invalidLoc()
		return
	}

	if status.portVisited && (!status.artistsSaved || !status.scientistsSaved) {
		updateMessage("Моряки непонимающе смотрят на вас" + pterm.Magenta("А где остальные?") + " — спрашивает капитан.\nВам нечего ответить. Вы зря потратили время.")
	}

	if !status.portVisited {
		updateMessage("Вы пришли в порт. Команда моряков вас уже ждала.\nРассказав им, что сейчас идут последнии часы Атлантиды, вы\nдаёте чёткое указание через час отплывать, не смотря ни на что.\nОдин из моряков просит у вас разрешения взять на корабль свинью.")
		options := []string{"Разрешить.", "Разрешить всем", "Запретить."}
		result, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Как поступить?").WithOptions(options).Show()
		switch result {
			case "Разрешить.":
				updateMessage("Вы разрешили ему сделать это. Несколько человек попросили взять и\nсвоих родственников, но вы отказали им. Некоторые остались недовольны.")
			case "Разрешить всем":
				updateMessage("Вы разрешили ему сделать это. Несколько человек попросили взять и\nсвоих родственников, вы им так же позволили сделать это.\nНо народу прошло много, а место на кораблях не безгранично.")
				status.sailorFamilySaved = true
			case "Запретить.":
				updateMessage("Вы отказали ему. Морякам это не понравилось, но они отнеслись с пониманием.")
		}
		if status.familySaved || status.scientistsSaved || status.artistsSaved {
			nextMessage()
		}
	}

	if status.familySaved {
		updateMessage(pterm.LightMagenta("Вы привели ребёнка, о болезни которого ничего не знаете!\nПожалуйста, подумайте о своём решение.") + " — сказал капитан.\nОн крайне недоволен, но не в праве вам перечить.\nВы подвергли всю команду этого судна опасности, ибо никому не ведома болезнь вашей дочери.")
		if status.scientistsSaved || status.artistsSaved {
			nextMessage()
		}
	}

	if status.artistsSaved || status.scientistsSaved {
		updateMessage("Вы готовы к отплытию?")
		options := []string{"Да", "Нет"}
		result, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Вы готовы к отплытию?").WithOptions(options).Show()
		if result == "Да" {
		endGame()
		} else {
			updateMessage("Нужно закончить дела...")
		}
	}
	status.portVisited = true
	resetLoc()
	in.port = true

}

func squareLoc() {
	updateMessage("Придя на центральную площадь города, вы встали на постамент и началить\nкричать, что есть мочи о надвигающейся опасности и короблях,\nчто стоят в порту. Народ повергла паника. Люди бегут к короблям, словно\nстадо испуганых овец. В порту началась бойня. Люди захватывают корабли,\nно их отбивают другие. Моряки перебиты. Надежды больше нет...")
	reset()
}
func artistsLoc() {
	if in.artists {
		invalidLoc()
		return
	}
	if status.artistsSaved {
		updateMessage("Вы же только что тут были, но почему-то опять пришли в уже пустое место.\nВы зря потеряли время.")
	}

	updateMessage("Вы считаете, что писатели, художники, поэты - люди искусств важны\nдля нового мира. Вы отправляетесь в гильдию искусств, чтобы предупредить\nлюдей о надвигающейся опасности и предложить спастись.\nЧлены гильдии вам благодарны, они идут на корабль.\nА время всё идёт... Нужно поторапливаться.")
	status.artistsSaved = true

	resetLoc()
	in.artists = true
}
func scientistsLoc() {
	if in.scientists {
		invalidLoc()
		return
	}

	if status.scientistsSaved {
		updateMessage("Учёные уже отправились на корабль. Вы зря сюда пришли.")
		in.scientists = true
		return
	}

	updateMessage("Вы считаете, что учёные крайне важны для создания нового общества.\nВы отправляетесь в резеденцию научного общества.\nОказывается там учёные неустанно трудятся над способом продлить жизнь Атлантиде.\nВы им рассказали о кораблях, на что они вам предложили ипользовать их\nизобретение, но они не уверены в его действенности.")

	options := []string{"Настоять на том, чтобы учёные отправились на корабль", "Использовать изобретение учёных", "Уйти"}
	result, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Что вы будете делать?").WithOptions(options).Show()
	switch result {
		case "Настоять на том, чтобы учёные отправились на корабль":
			updateMessage("Вы настаиваете, чтобы учёные всё бросили и спасались.\nОни не осмелились вас перечить, но считают, что вы обрекли королевство на гибель.")
			status.scientistsSaved = true
		case "Использовать изобретение учёных":
			updateMessage("Вы спросили учёных об изобретении. Они сказали, что его можно использовать\nтолько в центре острова и только в момент начала потопа. Выслушав все\nинструкции учёных, вы отправляетесь к центру острова.")
			nextMessage()
			if rand.Intn(3) == 0 {
				updateMessage("Устройство сработало! Буря тихнет! Мир спасён!")
				pterm.DefaultInteractiveSelect.WithDefaultText("Игра окончена").WithOptions([]string{"Продолжить"}).Show()
				reset()
				return
			} else {
				updateMessage("Устройство вырвалось у вас из рук, оно начало парить и ярко светиться,\nа через несколько секунд упало. Теперь это просто безделушка.\nВы обречены... Атлантида обречена...")
				reset()
				return
			}
		case "Уйти":
			updateMessage("Вы покидаете ученых, оставляя их работать над своим изобретением.\nВозможно это сработает, но у вас есть более важные дела.")
			return
	}

	resetLoc()
	in.scientists = true
}
func castleLoc() {
	if in.castle {
		invalidLoc()
		return
	}

	if status.artistsSaved || status.scientistsSaved {
		updateMessage("Вы пришли во дворец, чтобы попрощаться с королём. Вы больше никогда\nне увидитесь. Его лицо не выражало никаких эмоций, но вы понимали,\nчто его сейчас одолевала неумолимая горечь, страх и отчаяние.\nВы сказали, что набрали людей на корабль. Король ни сказал ни слова. Вам пора идти.")
		status.castleVisited = true
	} else {
		updateMessage(pterm.Magenta("Зачем ты пришёл сюда, у нас очень мало времени, поторопись!") + " — сказал король.\nВы зря потеряли время.")
	}

	resetLoc()
	in.castle = true
}
func homeLoc() {
	if in.home {
		invalidLoc()
		return
	}
	if status.familySaved {
		updateMessage("Ваша семья уже на борту, вы зря потрали время.")
	}
	if status.homeVisited {
		updateMessage("Ваша жена уже всё решила. Вам не переубедить её. Вы зря теряете время.")
		in.home = true
		return
	}

	updateMessage("Вы просто не могли оставить свою семью в неизвестности и поднимающейся панике.\nВторопях вы приходите домой. Жена рассказывает, что вашей дочери стало плохо.")
	options := []string{"Рассказать о катастрофе.", "Осмотреть дочь."}
	result := "" // Choice result
	doterChecked := false // Doter check status
	for result != "Рассказать о катастрофе." {
		result, _ = pterm.DefaultInteractiveSelect.WithDefaultText("Что вы будете делать?").WithOptions(options).Show()
		switch result {
			case "Рассказать о катастрофе.":
				updateMessage("Вы рассказываете жене о надвигающейся беде. Её ошеломило данное известие.")
			case "Осмотреть дочь.":
				updateMessage("Вы заходите в комнату дочери.\nОна вся красная и в поту. Ваша дочь серьёзно больна.")
				doterChecked = true
		}
	}
	options = []string{"Отправить семью на корабль."}
	if !doterChecked {
		options = append(options, "Осмотреть дочь.")
	}
	for result != "Отправить семью на корабль." {
		result, _ = pterm.DefaultInteractiveSelect.WithDefaultText("Время не ждёт").WithOptions(options).Show()
		switch result {
			case "Отправить семью на корабль.":
				updateMessage("Вы говорите семье, что они могут спастись.")
			case "Осмотреть дочь.":
				updateMessage("Вы заходите в комнату дочери.\nОна вся красная и в поту. Ваша дочь серьёзно больна.")
		}
	}
	nextMessage()
	if rand.Intn(3) == 0 {
		updateMessage("Ваша жена понимает, что они будет лишь обузой для остальных, а неизвестная\nболезнь дочери может распространиться на остальных. Они решают остаться.")
	} else {
		updateMessage("Ваша семья собирает вещи и отправляется на корабль.")
		status.familySaved = true
	}
	status.homeVisited = true
}
func endGame() {
	if !status.portVisited {
		updateMessage("Осталось уже совсем немного времени.  Небо затягивают тучи...")
	} else {
		updateMessage("Вы отправляетесь в порт. Осталось уже совсем немного времени.\nНебо затягивают тучи...")
	}
	nextMessage()

	if status.artistsSaved && status.scientistsSaved && status.sailorFamilySaved {
		updateMessage("Вы позвали слишком много народу. Пока люди не начали буйствовать вам нужно решать.")
		options := []string{"Оставить часть семей и учёных.", "Оставить часть семей и художников.", "Оставить нескольких художников и учёных.", "Оставить часть учёных и всех художников."}
		result, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Кто взойдёт на корабли?").WithOptions(options).Show()
		switch result {
			case "Оставить часть семей и учёных.":
				updateMessage("Вы решили сдержать обещание морякам и оставить хоть часть их семей,\nно так же нужно было взять с собой и учёных, нельзя было терять выдающиеся\nумы и научные достижения. К сожалению художникам места не досталось.\nНекоторые вышли по собственной воли, понимая, что так нужно, некоторые\nушли в знак протеста. Теперь на некоторых короблях неполная команда моряков,\nхудожники вот-вот взорвутся от яркости, " + pterm.Red("вы подставили их!") + "\nНо ничего уже не поделаешь...")

			case "Оставить часть семей и художников.":
				updateMessage("Вы решили сдержать обещание морякам и оставить хоть часть их семей,\nно так же нужно было взять с собой и учёных, нельзя было терять выдающиеся\nумы и научные достижения. К сожалению художникам места не досталось.\nНекоторые вышли по собственной воли, понимая, что так нужно, некоторые\nушли в знак протеста. Теперь на некоторых короблях неполная команда моряков,\nучёные вот-вот взорвутся от яркости, " + pterm.Red("вы подставили их!") + "\nНо ничего уже не поделаешь...")

			case "Оставить нескольких художников и учёных.", "Оставить часть учёных и всех художников.":
				updateMessage("Вы нарушили слово, данное морякам. Они подняли бунт. Вас, как и многих\nдругих выбросили за борт. Остался всего один корабль, и то, с неполной\nкомандой, но времени уже не было, буря началась. Почти все жители\nАтлантиды погибли. Вы были среди них. Ваше приключение закончилось.")
				reset()
				return
		}
	} else if status.sailorFamilySaved && (status.artistsSaved || status.scientistsSaved) {
		updateMessage("К сожалению не все влезают на корабль, нескольких людей пришлось оставить.\nНекоторые вышли по собственной воли, понимая, что так нужно, некоторые\nушли в знак протеста. Теперь на некоторых короблях не полная команда моряков,\nно ничего уже не поделаешь...")
	} else if status.artistsSaved && status.scientistsSaved {
		updateMessage("К сожалению людей слишком много. Не все смогут  взойти на судна.\nНекоторые добровольно отказались, но людей всё равно слишком много.")
		options := []string{"Предпочесть учёных.", "Предпочесть людей искусств."}
		result, _ := pterm.DefaultInteractiveSelect.WithDefaultText("Вам решать...").WithOptions(options).Show()
		switch result {
			case "Предпочесть учёных.":
				updateMessage("Вы считаете, что учёные важнее. К сожалению приходится оставить часть\nхудожников в Атлантиде. Напряжение нарастает...")
			case "Предпочесть людей искусств.":
				updateMessage("Вы считаете, что искусство важнее науки. К сожалению часть учёных\nпришлось оставить в Атлантиде. Напряжение нарастает...")
		}
	}
	nextMessage()
	updateMessage("Корабли, наполненные людьми отправляется в далёкое и долгое плавание.\nАтлантида скрывается под толщами воды прямо на выших глазах.")
	nextMessage()

	if status.familySaved {
		if rand.Intn(2) < 1 {
			updateMessage("Через пару недель вашей дочери стало лучше. К счастью болезнь миновала.\nВы смогли спасти семью.")
		} else {
			updateMessage("Через несколько недель плавания люди на вашем корабле начали чувствовать\nнемогание. Ещё через неделю заболелби почти все... Ваш корабль погиб.")
			reset()
			return
		}
	} else {
		updateMessage("Вы думаете о своей семье. Что с ними? Сожалеют ли они о своём выборе?")
	}
	nextMessage()
	updateMessage("Спустя пару месяцев все корабли достигли суши. Это место станет вашим\nновым домом. Возможно и оно когда-нибудь падёт, как это было с Атлантидой,\nно будем надеяться к тому времени мы будет готовы.")
	pterm.DefaultInteractiveSelect.WithDefaultText("Игра окончена").WithOptions([]string{"Ок"}).Show()
	reset()
	return
}
// Game timer
func endOfLife() {
	if status.time > 66 {
		updateMessage("Тучи сгущаются... Небеса разверзлись громом. Вот оно - начало конца.\nКорабли отчаливают. Вам не спастись.")
		nextMessage()
		// Count saved people
		savedPeople := 0
		if status.sailorFamilySaved{
			savedPeople += 1
		}
		if status.artistsSaved{
			savedPeople += 1
		}
		if status.scientistsSaved {
			savedPeople += 1
		}
		switch savedPeople {
			case 0:
				updateMessage("Практически все жители погибли. На кораблях была лишь их команда.")
			case 1:
				updateMessage("Лишь немногие люди смогут выжить.\n Вы не в их числе.")
			case 2,3: 
				updateMessage("Вы не смогли спастись, но организовали других. У человечества ещё есть шанс.")
		}
		nextMessage()
		if status.familySaved {
			updateMessage("Хоть вы и не успели, но ваша семья вместе с другими на корабле.\nХотя бы у них есть шанс.")
		} else {
			updateMessage("Вы и ваша семья обречена на гибель.\nОстаётся надеяться, что люди на кораблях справятся и\nдадут начало новому миру.")
		}
		reset()
		return
	}
	status.time += 17
}

