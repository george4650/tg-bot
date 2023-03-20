package main

import (
	"fmt"
	"io/ioutil"
	"log"
	Repository "myapp/internal/repository"
	"strconv"
	"strings"

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

	u := tgbotapi.NewUpdate(1)
	u.Timeout = 60

	updChannel, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic("update channel error: " + err.Error())
	}

	for update := range updChannel {

		if update.CallbackQuery != nil {

			if strings.Contains(update.CallbackQuery.Data, "Дoбавить товар в корзину") {

				var (
					Product_id []string
					product_id string
				)

				//Поучаем id товара
				Product_id = strings.Split(update.CallbackQuery.Data, "Дoбавить товар в корзину ")

				for _, p := range Product_id {
					product_id += p
				}

				product, err := Repository.ReadOne(product_id)
				if err != nil {
					log.Println(err)
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Возникла ошибка, товар не добавлен в корзину, попробуйте позже..."))
					continue
				}

				err = Repository.AddToCart(*product)
				if err != nil {
					log.Println(err)
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error()))
					continue
				}

				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(`Товар - "%s" добавлен в корзину`, product.Product_Name)))
			}

			if strings.Contains(update.CallbackQuery.Data, "Убрать товар из корзины с id = ") {

				var (
					Product_id []string
					product_id string
				)

				//Поучаем id товара
				Product_id = strings.Split(update.CallbackQuery.Data, "Убрать товар из корзины с id = ")

				for _, p := range Product_id {
					product_id += p
				}

				err := Repository.DeleteFromCart(product_id)
				if err != nil {
					log.Println(err)
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Возникла ошибка, товар не был удалён из корзины, попробуйте позже..."))
					continue
				}
				//bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(`Товар - "%s" удалён из корзины`, p.Product_Name)))
			}


			if strings.Contains(update.CallbackQuery.Data, "Увеличить товар в корзине с id = ") {

				var (
					Product_id []string
					product_id string
				)

				//Поучаем id товара
				Product_id = strings.Split(update.CallbackQuery.Data, "Увеличить товар в корзине с id = ")

				for _, p := range Product_id {
					product_id += p
				}

				err := Repository.IncrementKoll(product_id)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			if strings.Contains(update.CallbackQuery.Data, "Уменьшить товар в корзине с id = ") {

				var (
					Product_id []string
					product_id string
				)

				//Поучаем id товара
				Product_id = strings.Split(update.CallbackQuery.Data, "Уменьшить товар в корзине с id = ")

				for _, p := range Product_id {
					product_id += p
				}

				err := Repository.DeincrementKoll(product_id)
				if err != nil {
					log.Println(err)
					continue
				}
			}



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

				//Меню
				if update.Message.Text == mainMenu.Keyboard[0][0].Text {

					products, err := Repository.GetAllProducts()
					if err != nil {
						log.Println(err)
						continue
					}

					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Наше меню:")
					bot.Send(msgConfig)

					for _, p := range products {

						photoBytes, err := ioutil.ReadFile(p.Product_Image)
						if err != nil {
							log.Println(err)
						}

						photoFileBytes := tgbotapi.FileBytes{
							Name:  "picture",
							Bytes: photoBytes,
						}

						message := tgbotapi.NewPhotoUpload(int64(update.Message.Chat.ID), photoFileBytes)
						bot.Send(message)

						response := fmt.Sprintf("%s - %d руб\n", p.Product_Name, p.Product_Price)

						msgConfig = tgbotapi.NewMessage(update.Message.Chat.ID, response)

						msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("В корзину", fmt.Sprintf("Дoбавить товар в корзину %d", p.Product_Id)),
							),
						)
						bot.Send(msgConfig)
					}
				}

				//Мои заказы
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
						response := fmt.Sprintf("Заказ №%d\n %s - %d руб\n",
							key, o.Product_Name, o.Product_Price)

						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, response)

						bot.Send(msgConfig)
					}
				}

				//Корзина
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

					for _, o := range cart {

						photoBytes, err := ioutil.ReadFile(o.Product_Image)
						if err != nil {
							log.Println(err)
						}

						photoFileBytes := tgbotapi.FileBytes{
							Name:  "picture",
							Bytes: photoBytes,
						}

						message := tgbotapi.NewPhotoUpload(int64(update.Message.Chat.ID), photoFileBytes)
						bot.Send(message)

						response := fmt.Sprintf("%s - %d руб\n",
							o.Product_Name, o.Product_Price)

						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, response)

						msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("▼", fmt.Sprintf("Уменьшить товар в корзине с id = %d", o.Product_Id)),
								tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(o.Product_Koll), "noData"),
								tgbotapi.NewInlineKeyboardButtonData("▲", fmt.Sprintf("Увеличить товар в корзине с id = %d", o.Product_Id)),
							),
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Удалить", fmt.Sprintf("Убрать товар из корзины с id = %d", o.Product_Id)),
							),
						)

						 bot.Send(msgConfig)
					}
			
				}

			}

		}
	}

	bot.StopReceivingUpdates()
}
