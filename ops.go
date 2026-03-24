// Package mcrpc provides operator management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
	"github.com/fi-xz/mcrpc/internal/types"
)

// Operator is an alias for types.Operator, representing a server operator.
type Operator = types.Operator

// GetOperators retrieves the current list of operators.
func (c *MCRPCClient) GetOperators(context context.Context) ([]Operator, error) {
	var operators []Operator
	err := c.JSONRPCConn.Call(context, protocol.MethodOperatorsGet, nil, &operators)
	return operators, err
}

// SetOperators sets the operator list to the specified list, replacing the existing list.
func (c *MCRPCClient) SetOperators(context context.Context, operators []Operator) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.SetOperatorParams{Operators: operators}
	err := c.JSONRPCConn.Call(context, protocol.MethodOperatorsSet, params, &updatedOperators)
	return updatedOperators, err
}

// AddOperators adds the specified players as operators.
func (c *MCRPCClient) AddOperators(context context.Context, add []Operator) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.AddOperatorParams{OperatorAdd: add}
	err := c.JSONRPCConn.Call(context, protocol.MethodOperatorsAdd, params, &updatedOperators)
	return updatedOperators, err
}

// RemoveOperators removes the specified players from the operator list.
func (c *MCRPCClient) RemoveOperators(context context.Context, remove []Player) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.RemoveOperatorParams{OperatorRemove: remove}
	err := c.JSONRPCConn.Call(context, protocol.MethodOperatorsRemove, params, &updatedOperators)
	return updatedOperators, err
}

// ClearOperators removes all players from the operator list.
func (c *MCRPCClient) ClearOperators(context context.Context) ([]Operator, error) {
	var updatedOperators []Operator
	err := c.JSONRPCConn.Call(context, protocol.MethodOperatorsClear, nil, &updatedOperators)
	return updatedOperators, err
}
