package main

type nike struct {

}

func(a *nike) makeShoe() iShoe {
	return &nikeShoe {
		shoe: shoe {
			logo:"nike",
			size: 14,
		},
	}
}

func(a *nike) makeShirt() iShirt {
	return &nikeShirt {
		shirt: shirt {
			logo:"nike",
			size: 14,
		},
	}
}



