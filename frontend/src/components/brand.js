import React, { useState } from 'react';
import Row from 'react-bootstrap/Row'
import Form from 'react-bootstrap/Form'
import Card from 'react-bootstrap/Card'
import CardGroup from 'react-bootstrap/CardGroup'

function Brand(props) {
    let cards = props.data.map(el => (
        <CardGroup>
            <Card key={el.id}>
                <Card.Img variant="top" height={200} width={300} src={el.image_url} />
                <Card.Body>
                    <Card.Title>{el.name}</Card.Title>
                </Card.Body>
                <Card.Footer>
                    <small className="text-muted">{el.slug}</small>
                </Card.Footer>
            </Card>
        </CardGroup>
    ))
    return (
        <Row>
            {cards}
        </Row>
    );
}

export default Brand