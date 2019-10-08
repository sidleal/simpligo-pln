package senter

import (
	"log"
	"testing"
)

func TestSenter(t *testing.T) {

	// str := "1: A leishmaniose visceral, também chamada de calazar, é uma doença grave que atinge tanto humanos quanto animais. A doença, provocada pelo protozoário Leishmania infantum e transmitida pelo mosquito-palha Lutzomyia longipalpis, é responsável pelo sacrifício de inúmeros cães no Brasil.\n2: Braços galácticos são compostos por estrelas, gases e nuvens de poeira, que obstruem a visão que poderíamos ter do centro da galáxia. Esses componentes emitem diferentes tipos de radiação que, ao serem captados e interpretados, oferecem informações acerca de suas posições e velocidades. Tais objetos são chamados traçadores."
	// str := "A doença, provocada pelo protozoário Leishmania infantum e transmitida pelo mosquito-palha Lutzomyia longipalpis, é responsável pelo sacrifício de inúmeros cães no Brasil.\n2: Braços galácticos são compostos por estrelas, gases e nuvens de poeira, que obstruem a visão que poderíamos ter do centro da galáxia. Esses componentes emitem diferentes tipos de radiação que, ao serem captados e interpretados, oferecem informações acerca de suas posições e velocidades. Tais objetos são chamados traçadores."
	str := "1: A leishmaniose visceral, também chamada de calazar, é uma doença grave que atinge tanto humanos quanto animais. A doença, provocada pelo protozoário Leishmania infantum e transmitida pelo mosquito-palha Lutzomyia longipalpis é responsável pelo sacrifício de inúmeros cães no Brasil.\n2: Braços galácticos são compostos por estrelas, gases e nuvens de poeira, que obstruem a visão que poderíamos ter do centro da galáxia. Esses componentes emitem diferentes tipos de radiação que, ao serem captados e interpretados, oferecem informações acerca de suas posições e velocidades. Tais objetos são chamados traçadores."

	parsed := ParseText(str)

	for _, p := range parsed.Paragraphs {
		for _, s := range p.Sentences {
			for _, t := range s.Tokens {
				log.Println(t)
			}
		}
	}

}
