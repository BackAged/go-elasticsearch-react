import React, { useState } from 'react'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Form from 'react-bootstrap/Form'
import Brand from './brand'
import BrandApi from '../api/brand'

function Search() {
    const [state, setState] = useState({
        type: "brands",
        data: [],
    });

    const handleInputChange = (event) => {
        const cf = (data) => {
            console.log(data)
            setState((prevProps) => ({
                ...prevProps,
                "data": data,
            }))
            console.log(state)
        }
        BrandApi.SearchAsYouType(event.target.value, cf)
    };

    const handleRadioChange = (event) => {
        setState((prevProps) => ({
            ...prevProps,
            [event.target.name]: event.target.value
        }));
    }

    const handleSubmit = (event) => {
        event.preventDefault();
    };

    const isBrandSet = () => state.type === "brands"
    const isShopSet = () => state.type === "shops"
    const isProductSet = () => state.type === "products"

    return (
        <Row>
            <Row>
                <Form>
                    <Form.Row>
                        <Form.Check onChange={handleRadioChange} checked={isBrandSet()} inline name="type"
                            label="brands" value="brands" type="radio" id={`inline-brands`} />
                        <Form.Check onChange={handleRadioChange} checked={isShopSet()} inline name="type"
                            label="shops" value="shops" type="radio" id={`inline-shops`} />
                        <Form.Check onChange={handleRadioChange} checked={isProductSet()} inline name="type"
                            label="products" value="products" type="radio" id={`inline-products`} />
                    </Form.Row>
                    <Form.Row>
                        <Col>
                            <Form.Control onChange={handleInputChange} name="term" placeholder="search" />
                        </Col>
                    </Form.Row>
                </Form>
            </Row>
            <Row>
                {isBrandSet ? <Brand data={state.data} /> : "type something"}
            </Row>
        </Row>
    );
}

export default Search