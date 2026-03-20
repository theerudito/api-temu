import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";

const navLinks: { title: string; href: string }[] = [
  {
    title: "Pedidos Actuales",
    href: "/orders",
  },
{
    title: "Compradores",
    href: "/buyers",
  },
  {
    title: "Tabla de Precios",
    href: "/table",
  },

];

export function MainNavbar() {
  return (
    <NavigationMenu>
      <NavigationMenuList>
        {navLinks.map((item) => (
          <NavigationMenuItem key={item.title}>
            <NavigationMenuLink href={item.href} className={navigationMenuTriggerStyle()}>
              {item.title}
            </NavigationMenuLink>
          </NavigationMenuItem>
        ))}
      </NavigationMenuList>
    </NavigationMenu>
  );
}
