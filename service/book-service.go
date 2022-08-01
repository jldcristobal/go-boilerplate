package service

import (
	"fmt"
	"log"

	"github.com/breach-simulator/dto"
	"github.com/breach-simulator/entity"
	"github.com/breach-simulator/repository"
	"github.com/mashingan/smapping"
)

type BookService interface {
	Insert(book dto.BookCreateDTO) entity.Book
	Update(book dto.BookUpdateDTO) entity.Book
	Delete(book entity.Book)
	All() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (service *bookService) Insert(book dto.BookCreateDTO) entity.Book {
	bookToInsert := entity.Book{}
	err := smapping.FillStruct(&bookToInsert, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed to map %v", err)
	}
	res := service.bookRepository.InsertBook(bookToInsert)
	return res
}

func (service *bookService) Update(book dto.BookUpdateDTO) entity.Book {
	bookToUpdate := entity.Book{}
	err := smapping.FillStruct(&bookToUpdate, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed to map %v", err)
	}
	res := service.bookRepository.UpdateBook(bookToUpdate)
	return res
}

func (service *bookService) Delete(book entity.Book) {
	service.bookRepository.DeleteBook(book)
}

func (service *bookService) All() []entity.Book {
	return service.bookRepository.AllBooks()
}

func (service *bookService) FindByID(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	book := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", book.UserID)
	return userID == id
}
