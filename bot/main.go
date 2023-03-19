package main

import (
	"fmt"
	"log"
	Repository "myapp/internal/Repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const tgbotapiKey = "6212646389:AAEM-Y_FvE_S2-1nOU4lq8hWb20yzkgfEHU"

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📕 Меню"),
		tgbotapi.NewKeyboardButton("🗓 Мои заказы"),
		tgbotapi.NewKeyboardButton("🗒 Корзина"),
	),
)

func main() {

	err := Repository.OpenTable()
	if err != nil {
		panic("Connect to db error: " + err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(tgbotapiKey)
	if err != nil {
		panic("bot init error: " + err.Error())
	}

	botUser, err := bot.GetMe()
	if err != nil {
		panic("bot getme error: " + err.Error())
	}

	fmt.Printf("Авторизация прошла успешно! Запущен бот: %s\n", botUser.FirstName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updChannel, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic("update channel error: " + err.Error())
	}

	for update := range updChannel {

		if update.CallbackQuery != nil {

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))

		}

		if update.Message != nil {

			if update.Message.IsCommand() {
				Command := update.Message.Command()
				if Command == "start" {
					msg := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						fmt.Sprintf("Вы начали работу с ботом: %s", botUser.FirstName))
					msg.ReplyMarkup = mainMenu
					bot.Send(msg)
				}

			} else {

				if update.Message.Text == mainMenu.Keyboard[0][0].Text {

					products, err := Repository.GetAllProducts()
					if err != nil {
						log.Println(err)
						continue
					}

					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Наше меню:")
					bot.Send(msgConfig)

					for _, p := range products {
						response := fmt.Sprintf("%s - %d руб\n", p.Product_name, p.Product_price)

						msgConfig = tgbotapi.NewMessage(update.Message.Chat.ID, response)

						msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Просмотр", fmt.Sprintf("Просмотреть товар - %d", p.Product_id)),
								tgbotapi.NewInlineKeyboardButtonData("В корзину", fmt.Sprintf("Добавить товар в корзину - %d", p.Product_id)),
							),
						)
						bot.Send(msgConfig)
					}
				}
				if update.Message.Text == mainMenu.Keyboard[0][1].Text {

					orders, err := Repository.GetOrders(update.Message.From.UserName)
					if err != nil {
						log.Println(err)
						continue
					}
					if len(orders) == 0 {
						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас ещё не было заказов")
						bot.Send(msgConfig)
						continue
					}
					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "История ваших заказов:")
					bot.Send(msgConfig)

					for key, o := range orders {
						response := fmt.Sprintf("Заказ №%d\n %s - %s руб\n",
							key, o.Product_Name, o.Product_Price)

						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, response)

						bot.Send(msgConfig)
					}
				}
				if update.Message.Text == mainMenu.Keyboard[0][2].Text {

					if err != nil {
						log.Println(err)
						continue
					}

					cart := Repository.ReturnCart()
					if err != nil {
						log.Println(err)
						continue
					}

					if len(cart) == 0 {
						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваша корзина пуста")
						bot.Send(msgConfig)
						continue
					}
					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Корзина:")
					bot.Send(msgConfig)

					for key, o := range cart {
						response := fmt.Sprintf("%d) %s - %d руб\n",
							key, o.Product_name, o.Product_price)

						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, response)

						bot.Send(msgConfig)
					}
				} /* else {
					cs, ok := courseSignMap[update.Message.From.ID]
					if ok {
						if cs.State == finbot.StateEmail {
							cs.Email = update.Message.Text
							msgConfig := tgbotapi.NewMessage(
								update.Message.Chat.ID,
								"Введите телефон:")
							bot.Send(msgConfig)
							cs.State = 1
						} else if cs.State == finbot.StateTel {
							cs.Telephone = update.Message.Text
							cs.State = 2
							msgConfig := tgbotapi.NewMessage(
								update.Message.Chat.ID,
								"Введите course:")
							msgConfig.ReplyMarkup = courseMenu
							bot.Send(msgConfig)
						} else if cs.State == finbot.StateCourse {
							cs.Course = update.Message.Text
							msgConfig := tgbotapi.NewMessage(
								update.Message.Chat.ID,
								"ok!")
							msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
							bot.Send(msgConfig)
							delete(courseSignMap, update.Message.From.ID)
							//  post to site!
							err = post.SendPost(cs)
							if err != nil {
								fmt.Printf("send post error: %v\n", err)
							}
						}
						fmt.Printf("state: %+v\n", cs)
					} else {
						// other messages
						msgConfig := tgbotapi.NewMessage(
							update.Message.Chat.ID,
							"Не понятен ваш запрос")
						msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msgConfig)
					}
				}*/
			}
		}
	}

	bot.StopReceivingUpdates()
}
