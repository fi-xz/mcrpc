// Package mcrpc provides operator management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetOperators retrieves the current list of operators.
func (c *MCRPCClient) GetOperators(ctx context.Context) ([]Operator, error) {
	var operators []Operator
	err := c.JSONRPCConn.Call(ctx, protocol.MethodOperatorsGet, nil, &operators)
	return operators, err
}

// SetOperators sets the operator list to the specified list, replacing the existing list.
func (c *MCRPCClient) SetOperators(ctx context.Context, operators []Operator) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.SetOperatorParams{Operators: operators}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodOperatorsSet, params, &updatedOperators)
	return updatedOperators, err
}

// AddOperators adds the specified players as operators.
func (c *MCRPCClient) AddOperators(ctx context.Context, add []Operator) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.AddOperatorParams{OperatorAdd: add}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodOperatorsAdd, params, &updatedOperators)
	return updatedOperators, err
}

// RemoveOperators removes the specified players from the operator list.
func (c *MCRPCClient) RemoveOperators(ctx context.Context, remove []Player) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.RemoveOperatorParams{OperatorRemove: remove}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodOperatorsRemove, params, &updatedOperators)
	return updatedOperators, err
}

// ClearOperators removes all players from the operator list.
func (c *MCRPCClient) ClearOperators(ctx context.Context) ([]Operator, error) {
	var updatedOperators []Operator
	err := c.JSONRPCConn.Call(ctx, protocol.MethodOperatorsClear, nil, &updatedOperators)
	return updatedOperators, err
}
