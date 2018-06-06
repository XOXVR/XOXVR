const stripe = (data) => {
  const handler = StripeCheckout.configure({
    key: 'pk_test_0PO4bIauBnjhx8HPgWaUZYLZ',
    image: 'https://stripe.com/img/documentation/checkout/marketplace.png',
    locale: 'auto',
    token: (token) => {
      fetch('http://192.168.88.98:8080/api/v1/stripe/request', {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          token: token.id,
          email: token.email,
          amount: +data.cents,
          ethAddress: data.walletEthFiat
        }),
      })
        .then(res => res.json())
        .then((res) => {
          switch (res.status) {
            case 200:
              alert('Success! Check your tokens.');
              break;
            case 400:
              alert(res.error);
              break;
            default:
              alert('Server error. Status: 500!');
          }
        })
        .catch((err) => {
          alert('Сan not connect to the server!');
          console.log('error: ', err);
        });
    }
  });

  return handler;
};

export { stripe as default };
