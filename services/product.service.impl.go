package services

import (
	"context"
	"time"

	"github.com/akunsecured/emezen_api/models"
	"github.com/akunsecured/emezen_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductServiceImpl struct {
	productCollection *mongo.Collection
	userService       UserService
	ctx               context.Context
}

func NewProductService(productCollection *mongo.Collection, userService UserService, ctx context.Context) ProductService {
	return &ProductServiceImpl{
		productCollection: productCollection,
		userService:       userService,
		ctx:               ctx,
	}
}

func (p *ProductServiceImpl) AddProduct(product *models.Product) (*string, error) {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	result, err := p.productCollection.InsertOne(p.ctx, product)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		var oidHex = oid.Hex()
		return &oidHex, nil
	}
	return nil, utils.ErrInsertedIDIsNotObjectID
}

func (p *ProductServiceImpl) GetProduct(productId *string) (*models.Product, error) {
	var product *models.Product
	objID, err := primitive.ObjectIDFromHex(*productId)
	if err != nil {
		return nil, err
	}
	query := bson.D{bson.E{Key: "_id", Value: objID}}
	err = p.productCollection.FindOne(p.ctx, query).Decode(&product)
	return product, err
}

func (p *ProductServiceImpl) GetAllProducts() ([]*models.Product, error) {
	var products []*models.Product

	cur, err := p.productCollection.Find(p.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(p.ctx) {
		//Create a value into which the single document can be decoded
		var product *models.Product
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	cur.Close(p.ctx)

	return products, err
}

func (p *ProductServiceImpl) GetAllProductsOfUser(userId *string) ([]*models.Product, error) {
	products, err := p.GetAllProducts()
	if err != nil {
		return nil, err
	}

	var result []*models.Product

	for _, product := range products {
		if product.SellerID == *userId {
			result = append(result, product)
		}
	}

	return result, err
}

func (p *ProductServiceImpl) UpdateProduct(product *models.Product) error {
	filter := bson.D{bson.E{Key: "_id", Value: product.ID}}

	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "seller_id", Value: product.SellerID},
		bson.E{Key: "name", Value: product.Name},
		bson.E{Key: "price", Value: product.Price},
		bson.E{Key: "images", Value: product.Images},
		bson.E{Key: "details", Value: product.Details},
		bson.E{Key: "quantity", Value: product.Quantity},
		bson.E{Key: "category", Value: product.Category},
		bson.E{Key: "updated_at", Value: time.Now()},
	}}}
	result, err := p.productCollection.UpdateOne(p.ctx, filter, update)
	if result.MatchedCount != 1 {
		return utils.ErrNotExists
	}
	return err
}

func (p *ProductServiceImpl) DeleteProduct(productId *string) error {
	objID, err := primitive.ObjectIDFromHex(*productId)
	if err != nil {
		return err
	}

	filter := bson.D{bson.E{Key: "_id", Value: objID}}
	result, _ := p.productCollection.DeleteOne(p.ctx, filter)
	if result.DeletedCount != 1 {
		return utils.ErrNoMatchedDocumentFoundForDelete
	}

	return nil
}
