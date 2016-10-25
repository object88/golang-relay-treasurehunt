package data

import (
	"github.com/graphql-go/graphql"
	"github.com/object88/relay"
	"golang.org/x/net/context"
)

var gameType *graphql.Object
var hidingSpotType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions
var hidingSpotConnection *relay.GraphQLConnectionDefinitions

// Schema is our published GraphQL representation of objects and mutations
var Schema graphql.Schema

func init() {
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ct context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			if resolvedID.Type == "Game" {
				return GetGame(), nil
			}
			if resolvedID.Type == "HidingSpot" {
				id := resolvedID.ID
				return GetGame().GetHidingSpot(id), nil
			}
			return nil, nil
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			switch p.Value.(type) {
			case *Game:
				return gameType
			case *HidingSpot:
				return hidingSpotType
			}
			return nil
		},
	})

	hidingSpotType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "HidingSpot",
		Description: "A place where you might find treasure",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("HidingSpot", nil),
			"hasBeenChecked": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "True if this spot has already been checked for treasure",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					hidingSpot := p.Source.(*HidingSpot)
					return hidingSpot.HasBeenChecked, nil
				},
			},
			"hasTreasure": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "True if this hiding spot holds treasure",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					hidingSpot := p.Source.(*HidingSpot)
					if hidingSpot.HasBeenChecked {
						return hidingSpot.HasTreasure, nil
					}
					return nil, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	hidingSpotConnection = relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "HidingSpotConnection",
		NodeType: hidingSpotType,
	})

	gameType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Game",
		Description: "A treasure search game",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Game", nil),
			"hidingSpots": &graphql.Field{
				Type:        hidingSpotConnection.ConnectionType,
				Description: "Places where treasure might be hidden",
				Args:        relay.ConnectionArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					arg := relay.NewConnectionArguments(p.Args)
					hidingSpots := GetGame().GetHidingSpots()
					interfaceSlice := hidingSpotToInterfaceSlice(hidingSpots...)
					return relay.ConnectionFromArray(interfaceSlice, arg), nil
				},
			},
			"turnsRemaining": &graphql.Field{
				Type:        graphql.Int,
				Description: "The number of turns a player has left to find the treasure",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetGame().GetTurnsRemaining(), nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	checkHidingSpotForTreasureMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "CheckHidingSpotForTreasure",
		InputFields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		OutputFields: graphql.Fields{
			"hidingSpot": {
				Type: hidingSpotType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						return GetGame().GetHidingSpot(payload["id"].(string)), nil
					}
					return nil, nil
				},
			},
			"game": {
				Type: gameType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetGame(), nil
				},
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			id := inputMap["id"].(string)
			localHidingSpotID := relay.FromGlobalID(id).ID
			GetGame().CheckHidingSpotForTreasure(localHidingSpotID)
			return map[string]interface{}{
				"id": localHidingSpotID,
			}, nil
		},
	})

	/**
	 * This is the type that will be the root of our mutations,
	 * and the entry point into performing writes in our schema.
	 */
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"checkHidingSpotForTreasure": checkHidingSpotForTreasureMutation,
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": nodeDefinitions.NodeField,
			"game": &graphql.Field{
				Type: gameType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetGame(), nil
				},
			},
		},
	})

	/**
	 * Finally, we construct our schema (whose starting query type is the query
	 * type we defined above) and export it.
	 */
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Mutation: mutationType,
		Query:    queryType,
		Types:    []graphql.Type{gameType, hidingSpotType},
	})
	if err != nil {
		panic(err)
	}
}

func hidingSpotToInterfaceSlice(hidingSpots ...*HidingSpot) []interface{} {
	var interfaceSlice = make([]interface{}, len(hidingSpots))
	for i, d := range hidingSpots {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}
