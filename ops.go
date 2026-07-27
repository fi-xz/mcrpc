// Package mcrpc provides operator management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetOperators retrieves the current list of operators.
func (c *Client) GetOperators(ctx context.Context) ([]Operator, error) {
	var operators []Operator
	err := c.call(ctx, protocol.MethodOperatorsGet, nil, &operators)
	return operators, err
}

// SetOperators sets the operator list to the specified list, replacing the existing list.
func (c *Client) SetOperators(ctx context.Context, operators []Operator) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.SetOperatorParams{Operators: operators}
	err := c.call(ctx, protocol.MethodOperatorsSet, params, &updatedOperators)
	return updatedOperators, err
}

// AddOperators adds the specified players as operators.
func (c *Client) AddOperators(ctx context.Context, add []Operator) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.AddOperatorParams{OperatorAdd: add}
	err := c.call(ctx, protocol.MethodOperatorsAdd, params, &updatedOperators)
	return updatedOperators, err
}

// RemoveOperators removes the specified players from the operator list.
func (c *Client) RemoveOperators(ctx context.Context, remove []Player) ([]Operator, error) {
	var updatedOperators []Operator
	params := protocol.RemoveOperatorParams{OperatorRemove: remove}
	err := c.call(ctx, protocol.MethodOperatorsRemove, params, &updatedOperators)
	return updatedOperators, err
}

// ClearOperators removes all players from the operator list.
func (c *Client) ClearOperators(ctx context.Context) ([]Operator, error) {
	var updatedOperators []Operator
	err := c.call(ctx, protocol.MethodOperatorsClear, nil, &updatedOperators)
	return updatedOperators, err
}
